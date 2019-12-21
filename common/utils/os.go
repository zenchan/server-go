package utils

import (
	"os"
	"runtime"
	"strings"
)

var originOSArgs []string

func init() {
	originOSArgs = make([]string, len(os.Args))
	for i, arg := range os.Args {
		originOSArgs[i] = arg
	}
}

var procName string

// ProcessName returns process name
func ProcessName() string {
	if procName == "" {
		segs := strings.Split(originOSArgs[0], string(os.PathSeparator))
		switch runtime.GOOS {
		case "windows":
			procName = strings.Split(segs[len(segs)-1], ".")[0]
		case "linux":
			procName = segs[len(segs)-1]
		}
	}
	return procName
}
