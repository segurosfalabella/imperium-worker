#/bin/bash

go get -t ./...
GO_FILES=$(find . -iname '*.go' -type f | grep -v /godogs/ | tr "\n" " ") # All the .go files, excluding vendor/
go get github.com/golang/lint/golint                        # Linter
go get github.com/fzipp/gocyclo
go get github.com/DATA-DOG/godog/cmd/godog


test -z $(gofmt -s -l $GO_FILES)
go test -v -race ./...
go vet -v ./...
cd godogs
godog
cd ..
gocyclo -over 4 $(echo $GO_FILES)
golint -set_exit_status $(go list ./...)
$GOPATH/bin/goveralls -service=travis-ci
