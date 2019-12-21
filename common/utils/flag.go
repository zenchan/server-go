package utils

import (
	"flag"
	"log"
)

// flag parameters
var (
	ConfigFile string
)

// FlagParse parses the command-line flags
func FlagParse() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pname := ProcessName()
	flag.StringVar(&ConfigFile, "c", pname+".ini", "configuration file")
	flag.Parse()
}
