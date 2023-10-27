.PHONY: all build run gotool clean help

BINARY="go_web_app"

all: gotool build

build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${BINARY}

run:
	go run ./main.go --c conf/config.yaml

gotool:
	go fmt ./...
	go vet ./...

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

