package xlog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	debugLog = iota
	infoLog
	warnLog
	errorLog
	fatalLog
)

var logLevelName = []string{
	debugLog: "DEBUG",
	infoLog:  "INFO",
	warnLog:  "WARN",
	errorLog: "ERROR",
	fatalLog: "FATAL",
}

type loggingT struct {
	// skip is the number of goroutine stack frames to skip before call xlog to log
	skip        int32
	logLevel    int32
	logDir      string // log file's base path
	logDateDir  string
	logFile     string // log file's name
	logFullName string
	file        *os.File
	mu          sync.Mutex
	freeList    *sync.Pool
	now         time.Time
	rotateTime  int64 // next time to rotate file

	// options
	logStdout int32
}

var logging loggingT

// InitLogging initialize xlog
func InitLogging(logDir, logLevel string, opts ...Option) error {
	dDir := dateDir()
	dir := logDir + "/" + dDir
	if err := createLogDir(dir); err != nil {
		return err
	}

	logging.now = time.Now()
	hour, _, _ := logging.now.Clock()
	t, _ := time.ParseInLocation("20060102 15", dDir+" "+strconv.Itoa(hour+1), logging.now.Location())
	logging.rotateTime = t.Unix()
	logging.logDir = logDir
	logging.logDateDir = dDir
	logging.logFile = getFilename(&logging.now)
	logging.logFullName = dir + "/" + logging.logFile

	f, err := os.OpenFile(logging.logFullName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logging.file = f

	logging.logLevel = 0 // debugLog
	logging.setLevel(logLevel)

	logging.freeList = &sync.Pool{
		New: func() interface{} {
			return &buffer{}
		},
	}

	for _, opt := range opts {
		opt(&logging)
	}

	go logging.recoverFile()

	return nil
}

// recoverFile recovers log file whenever log file or directory been deleted.
func (l *loggingT) recoverFile() {
	tk := time.NewTicker(time.Second)
	for {
		select {
		case now := <-tk.C:
			l.mu.Lock()
			if now.Unix() >= l.rotateTime { // new log file
				l.now = now
				l.rotateFile()
			} else { // current log file
				info1, err := os.Stat(l.logFullName)
				info2, _ := l.file.Stat()
				if err != nil || !os.SameFile(info1, info2) { // log file been deleted
					l.file.Close()
					l.file = nil
					l.rotateFile()
				}
			}
			l.mu.Unlock()
		}
	}
}

func (l *loggingT) header(level int) *buffer {
	buf := l.freeList.Get().(*buffer)
	buf.Reset()
	skip := l.getSkip()
	pc, file, line, ok := runtime.Caller(3 + int(skip))
	var fn string
	if !ok {
		file = "???"
		line = 1
		fn = "main"
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
		fn = runtime.FuncForPC(pc).Name()
	}
	return l.formatHeader(buf, level, file, line, fn)
}

// formatHeader formats a log header using the provided file name, line number and function name.
func (l *loggingT) formatHeader(buf *buffer, level int, file string, line int, fn string) *buffer {
	l.now = time.Now()
	year, month, day := l.now.Date()
	hour, min, sec := l.now.Clock()
	_, offset := l.now.Zone()

	buf.nDigits(4, 0, year, '0')
	buf.tmp[4] = '-'
	buf.twoDigits(5, int(month))
	buf.tmp[7] = '-'
	buf.twoDigits(8, day)
	buf.tmp[10] = ' '
	buf.twoDigits(11, hour)
	buf.tmp[13] = ':'
	buf.twoDigits(14, min)
	buf.tmp[16] = ':'
	buf.twoDigits(17, sec)
	buf.tmp[19] = '.'
	buf.nDigits(6, 20, l.now.Nanosecond()/1000, '0')
	buf.tmp[26] = ' '
	zone := offset / 60
	if zone < 0 {
		buf.tmp[27] = '-'
	} else {
		buf.tmp[27] = '+'
	}
	buf.twoDigits(28, zone/60)
	buf.twoDigits(30, zone%60)
	buf.tmp[32] = '\t'
	n := len(logLevelName[level])
	copy(buf.tmp[33:], logLevelName[level][:])
	n = 33 + n
	buf.tmp[n] = '\t'
	n++
	n += buf.someDigits(n, pid)
	buf.tmp[n] = '\t'
	buf.Write(buf.tmp[:n+1])
	buf.WriteString(file)
	buf.tmp[0] = ':'
	n = 1 + buf.someDigits(1, line)
	buf.tmp[n] = '\t'
	buf.Write(buf.tmp[:n+1])
	if i := strings.LastIndex(fn, "/"); i > 0 {
		buf.WriteString(fn[i+1:])
	} else {
		buf.WriteString(fn)
	}
	buf.WriteByte('\t')

	return buf
}

type buffer struct {
	bytes.Buffer
	tmp [64]byte
}

const digits = "0123456789"

func (buf *buffer) twoDigits(i, d int) {
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}

func (buf *buffer) nDigits(n, i, d int, pad byte) {
	j := n - 1
	for ; j >= 0 && d > 0; j-- {
		buf.tmp[i+j] = digits[d%10]
		d /= 10
	}
	for ; j >= 0; j-- {
		buf.tmp[i+j] = pad
	}
}

func (buf *buffer) someDigits(i, d int) int {
	j := len(buf.tmp)
	for d > 0 {
		j--
		buf.tmp[j] = digits[d%10]
		d /= 10
	}
	return copy(buf.tmp[i:], buf.tmp[j:])
}

func (l *loggingT) print(level int, args ...interface{}) {
	lv := l.getLevel()
	if level < lv {
		return
	}

	buf := l.header(level)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(level, buf)
}

func (l *loggingT) printf(level int, format string, args ...interface{}) {
	lv := l.getLevel()
	if level < lv {
		return
	}

	buf := l.header(level)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(level, buf)
}

func (l *loggingT) output(level int, buf *buffer) {
	data := buf.Bytes()
	l.mu.Lock()
	if l.file == nil {
		os.Stdout.Write([]byte("ERROR: logging before InitLogging\n"))
		os.Stdout.Write(data)
		l.mu.Unlock()
		return
	}

	if err := l.rotateFile(); err != nil {
		os.Stdout.Write([]byte(fmt.Sprintf("ERROR: logging rotate file failed, %s\n", err.Error())))
		os.Stdout.Write(data)
		l.mu.Unlock()
		return
	}

	l.file.Write(data)
	if l.stdout() == 1 {
		os.Stdout.Write(data)
	}
	l.mu.Unlock()
	l.freeList.Put(buf)
}

// rotateFile rotates log file every hour.
func (l *loggingT) rotateFile() error {
	if l.now.Unix() >= l.rotateTime {
		l.logDateDir = dateDir()
		l.logFullName = l.logDir + "/" + l.logDateDir + "/" + getFilename(&l.now)
		l.file.Close()
		l.file = nil
		l.rotateTime += 3600 // next time to rotate file is an hour later
	}

	if l.file == nil {
		dir := l.logDir + "/" + l.logDateDir
		if _, err := os.Stat(dir); err != nil {
			if err := createLogDir(dir); err != nil {
				return err
			}
		}

		f, err := os.OpenFile(l.logFullName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		l.file = f
	}

	return nil
}

func (l *loggingT) getLevel() int {
	return int(atomic.LoadInt32(&l.logLevel))
}

func (l *loggingT) levelString() string {
	lv := l.getLevel()
	s := logLevelName[lv]
	return s
}

func (l *loggingT) setLevel(level string) {
	level = strings.ToUpper(level)
	for i := range logLevelName {
		if logLevelName[i] == level {
			atomic.StoreInt32(&l.logLevel, int32(i))
			break
		}
	}
}

func (l *loggingT) getSkip() int {
	return int(atomic.LoadInt32(&l.skip))
}

func (l *loggingT) setSkip(skip int) {
	atomic.StoreInt32(&l.skip, int32(skip))
}

func (l *loggingT) stdout() int32 {
	return atomic.LoadInt32(&l.logStdout)
}

// LevelString returns log level string
func LevelString() string {
	return logging.levelString()
}

// SetLevel sets log level
func SetLevel(level string) {
	logging.setLevel(level)
}

// SetSkip sets the number of goroutine stack frames to skip
func SetSkip(skip int) {
	logging.setSkip(skip)
}

// SetStdout whether output to stdout when write to log
func SetStdout(b bool) {
	if b {
		atomic.StoreInt32(&logging.logStdout, 1)
	} else {
		atomic.StoreInt32(&logging.logStdout, 0)
	}
}

// Debug logs to the DEBUG log.
// A newline is appended if missing.
func Debug(args ...interface{}) {
	logging.print(debugLog, args...)
}

// Debugln logs to the DEBUG log.
// A newline is appended if missing.
func Debugln(args ...interface{}) {
	logging.print(debugLog, args...)
}

// Debugf logs to the DEBUG log.
// A newline is appended if missing.
func Debugf(format string, args ...interface{}) {
	logging.printf(debugLog, format, args...)
}

// Info logs to the INFO log.
// A newline is appended if missing.
func Info(args ...interface{}) {
	logging.print(infoLog, args...)
}

// Infoln logs to the INFO log.
// A newline is appended if missing.
func Infoln(args ...interface{}) {
	logging.print(infoLog, args...)
}

// Infof logs to the INFO log.
// A newline is appended if missing.
func Infof(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}

// Warning logs to the WARNING log.
// A newline is appended if missing.
func Warning(args ...interface{}) {
	logging.print(warnLog, args...)
}

// Warningln logs to the WARNING log.
// A newline is appended if missing.
func Warningln(args ...interface{}) {
	logging.print(warnLog, args...)
}

// Warningf logs to the WARNING log.
// A newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	logging.printf(warnLog, format, args...)
}

// Error logs to the ERROR log.
// A newline is appended if missing.
func Error(args ...interface{}) {
	logging.print(errorLog, args...)
}

// Errorln logs to the ERROR log.
// A newline is appended if missing.
func Errorln(args ...interface{}) {
	logging.print(errorLog, args...)
}

// Errorf logs to the ERROR log.
// A newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	logging.printf(errorLog, format, args...)
}

// Fatal logs to the FATAL log.
// A newline is appended if missing.
func Fatal(args ...interface{}) {
	logging.print(fatalLog, args...)
}

// Fatalln logs to the FATAL log.
// A newline is appended if missing.
func Fatalln(args ...interface{}) {
	logging.print(fatalLog, args...)
}

// Fatalf logs to the FATAL log.
// A newline is appended if missing.
func Fatalf(format string, args ...interface{}) {
	logging.printf(fatalLog, format, args...)
}
