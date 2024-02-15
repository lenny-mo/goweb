.PHONY: all build run gotool clean help \
 	vegetaWithReport \
 	vegetaWithEncodeStorage	\
 	vegetaGenerateHtml

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

vegetaWithReport:
	vegeta -cpus 1 attack \
    -targets=./vegeta/targets.txt \
    -duration=30s \
    -connections=100 \
    -rate=500 \
    -format=http | \
    vegeta report -type=text -every=2s | tee vegetaTestWithReport.txt

vegetaWithEncodeStorage:
	vegeta -cpus 1 attack \
    -targets=./vegeta/targets.txt \
    -duration=30s \
    -connections=100 \
    -rate=500 \
    -format=http | \
    vegeta encode > vegetaTestWithEncode.json

`vegetaGenerateHtml`: vegetaWithEncodeStorage
	 vegeta plot vegetaTestWithEncode.json > plot.html
