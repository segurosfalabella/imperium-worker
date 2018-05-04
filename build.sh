#/bin/sh

export CGO_ENABLED=0
export GOOS=linux

go get -t ./...
go build -a -installsuffix cgo
