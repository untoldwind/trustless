VERSION ?= $(shell date -u +%Y%m%d.%H%M%S)

all: export GOPATH=${PWD}/../../../..
all: format
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-s -w -X github.com/untoldwind/scrypt.go/config.version=${VERSION}" -v -i -o bin/scrypt-go github.com/untoldwind/scrypt-go

format: export GOPATH=${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ./...

test: export GOPATH=${PWD}/../../../..
test:
	@echo "--> Running go test"
	@go test -v ./...

dep.install:
	@echo "-> dep install"
	@go get github.com/golang/dep/cmd/dep
	@go build -v -o bin/dep github.com/golang/dep/cmd/dep

dep.ensure: dep.install
	@bin/dep ensure
	@bin/dep prune
	