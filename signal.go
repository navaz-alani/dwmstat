package main

import (
	"io/ioutil"
	l "log"
	"net"
	"os"
)

const SIG_SOCK = "/tmp/" + EXE + "_sig"

type Signal struct {
	Module string
}

func listenSignals(ss chan<- Signal) {
	if err := os.RemoveAll(SIG_SOCK); err != nil {
		l.Fatal("socket remove error:", err)
	}

	sock, err := net.Listen("unix", SIG_SOCK)
	if err != nil {
		l.Fatal("socket listen error:", err)
	}
	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			log(ERROR, "failed to accept connection: %s", err)
		}
		go func(conn net.Conn) {
			if data, err := ioutil.ReadAll(conn); err != nil {
				log(ERROR, "failed to read from connection: %s", err)
			} else {
				ss <- Signal{string(data)}
			}
		}(conn)
	}
}
