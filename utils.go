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
	l.Printf("["+EXE+":"+kind+"] "+s, args...)
}
