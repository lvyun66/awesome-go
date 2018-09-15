FROM golang:latest

WORKDIR $GOPATH/src/github.com/lvyun66/awesome-go/netease

RUN go get github.com/lvyun66/awesome-go & \
go build .

ENTRYPOINT ./netease