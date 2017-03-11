PACKAGES=$(shell go list ./...)
VERSION = $(shell date -u +.%Y%m%d.%H%M%S)

all: export GOPATH=${PWD}/../../../..
all: format
	@echo "--> Running go build"
	@go build ./...

format: export GOPATH=${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ./...

test: export GOPATH=${PWD}/../../../..
test:
	@echo "--> Running tests"
	@go test -v . ./rest/...

godepssave:
	@echo "--> Godeps save"
	@go get github.com/tools/godep
	@go build -v -o bin/godep github.com/tools/godep
	@bin/godep save ./...
