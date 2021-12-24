EXE=dwmstat
GO_SRC=$(wildcard *.go)

.PHONY: install clean

${EXE}: ${GO_SRC}
	go build -o $@

install: ${EXE}
	cp -f $^ ${DESTDIR}/bin/

clean:
	rm -rf ${EXE}
