from golang:1.20.10-alpine3.18 as builder

env GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#当前的工作目录
workdir /build

copy go.mod .
copy go.sum .
run go mod download
#复制所有代码文件到当前工作目录
copy . .

run go build -o go_web_app  #编译go_web_app

## final stage 需要执行shell命令
from debian:buster
copy ./conf /conf
copy ./wait-for.sh /
copy ./mysql_init.sql /

#复制编译好的go_web_app到根目录
copy --from=builder /build/go_web_app /

run set -eux \
    && apt-get update \
    && apt-get install -y --no-install-recommends netcat \
    && chmod +x wait-for.sh;

#entrypoint ["/go_web_app", "-c", "/conf/config.yaml"]

