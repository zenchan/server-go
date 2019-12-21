package xlog

import (
	"os"
	"time"
)

var (
	pid = os.Getpid()
)

func createLogDir(dir string) error {
	var err error
	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
		}
	}
	return err
}

func dateDir() string {
	now := time.Now()
	year, month, day := now.Date()

	var buf [8]byte
	buf[3] = digits[year%10]
	year /= 10
	buf[2] = digits[year%10]
	year /= 10
	buf[1] = digits[year%10]
	year /= 10
	buf[0] = digits[year]

	mm := int(month)
	buf[5] = digits[mm%10]
	mm /= 10
	buf[4] = digits[mm]

	buf[7] = digits[day%10]
	day /= 10
	buf[6] = digits[day]

	return string(buf[:])
}

// getFilename gets log file name
func getFilename(now *time.Time) string {
	hour, _, _ := now.Clock()
	var buf [2]byte
	buf[1] = digits[hour%10]
	hour /= 10
	buf[0] = digits[hour]
	return string(buf[:]) + ".log"
}
