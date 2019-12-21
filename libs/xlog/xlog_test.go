package xlog

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

const testLogDir = "./testlog"

func initTestLog(level string) error {
	return InitLogging(testLogDir, level)
}

func cleanTestLog() {
	if logging.file != nil {
		logging.file.Close()
		logging.file = nil
	}
	os.RemoveAll(testLogDir)
}

func contains(level int, msg string, t *testing.T) bool {
	if logging.file != nil {
		logging.file.Close()
		logging.file = nil
	}

	buf, err := ioutil.ReadFile(logging.logFullName)
	if err != nil {
		t.Logf("read file %s failed, %s", logging.logFullName, err.Error())
		return false
	}

	s := logLevelName[level]
	cont := string(buf)
	if !strings.Contains(cont, s) {
		return false
	}
	if !strings.Contains(cont, msg) {
		return false
	}

	return true
}

var logTestData = []struct {
	level   int
	fn      func(args ...interface{})
	msg     string
	testMsg string
}{
	{debugLog, Debug, "debuglog", "Debug"},
	{infoLog, Info, "infolog", "Info"},
	{warnLog, Warning, "warninglog", "Warning"},
	{errorLog, Error, "errorlog", "Error"},
	{fatalLog, Fatal, "fatallog", "Fatal"},
}

func TestDebug(t *testing.T) {
	if err := initTestLog("debug"); err != nil {
		t.Fatalf("initTestLog failed, %s", err.Error())
	}

	for _, d := range logTestData {
		d.fn(d.msg)
	}

	for _, d := range logTestData {
		if !contains(d.level, d.msg, t) {
			t.Error(d.testMsg + " failed")
		}
	}

	cleanTestLog()
}

func TestInfo(t *testing.T) {
	if err := initTestLog("info"); err != nil {
		t.Fatalf("initTestLog failed, %s", err.Error())
	}

	for _, d := range logTestData {
		d.fn(d.msg)
	}

	for _, d := range logTestData {
		if d.level < infoLog {
			if contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		} else {
			if !contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		}
	}

	cleanTestLog()
}

func TestWarning(t *testing.T) {
	if err := initTestLog("warning"); err != nil {
		t.Fatalf("initTestLog failed, %s", err.Error())
	}

	for _, d := range logTestData {
		d.fn(d.msg)
	}

	for _, d := range logTestData {
		if d.level < warnLog {
			if contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		} else {
			if !contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		}
	}

	cleanTestLog()
}

func TestError(t *testing.T) {
	if err := initTestLog("error"); err != nil {
		t.Fatalf("initTestLog failed, %s", err.Error())
	}

	for _, d := range logTestData {
		d.fn(d.msg)
	}

	for _, d := range logTestData {
		if d.level < errorLog {
			if contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		} else {
			if !contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		}
	}

	cleanTestLog()
}

func TestFatal(t *testing.T) {
	if err := initTestLog("fatal"); err != nil {
		t.Fatalf("initTestLog failed, %s", err.Error())
	}

	for _, d := range logTestData {
		d.fn(d.msg)
	}

	for _, d := range logTestData {
		if d.level < fatalLog {
			if contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		} else {
			if !contains(d.level, d.msg, t) {
				t.Error(d.testMsg + " failed")
			}
		}
	}

	cleanTestLog()
}

func BenchmarkHeader(b *testing.B) {
	logging.freeList = &sync.Pool{
		New: func() interface{} {
			return &buffer{}
		},
	}

	for i := 0; i < b.N; i++ {
		buf := logging.header(0)
		logging.freeList.Put(buf)
	}
}

func BenchmarkLog(b *testing.B) {
	if err := initTestLog("debug"); err != nil {
		b.Fatalf("initTestLog failed, %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		Info("this is a test log")
	}

	cleanTestLog()
}
