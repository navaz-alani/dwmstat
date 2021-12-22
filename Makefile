EXE=dwmstat
GO_SRC=$(wildcard *.go)

.PHONY: install clean

${EXE}: ${GO_SRC}
	go build -o $@

install: ${EXE}
	# kill EXE if already running
	bash -c 'pkill $^ || true'
	cp $^ ../bin/
	# spawn a new instance of EXE
	${EXE} &>/dev/null  &

clean:
	rm -rf ${EXE}
