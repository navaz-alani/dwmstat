package main

import (
	l "log"
)

const (
	EXE = "dwmstat"

	ERROR = "ERRO"
	INFO  = "INFO"
	WARN  = "WARN"
)

func log(kind, s string, args ...interface{}) {
	// only log warnings and errors in production
	if *production {
		if kind == WARN || kind == ERROR {
			l.Printf("["+EXE+":"+kind+"] "+s, args...)
		}
	} else {
		l.Printf("["+EXE+":"+kind+"] "+s, args...)
	}
}
