package main

import (
	"flag"
	l "log"
)

// parse command line arguments and perform any checks
func init() {
	flag.Parse()

	if *SIG_SOCK == "" {
		l.Fatalln("socket unspecified!")
	}
}

func main() { statusBar.Run() }
