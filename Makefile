PKGS = ./api/... ./commands/... ./config/... ./daemon/... ./secrets/... ./store/... .

VERSION ?= $(shell date -u +%Y%m%d.%H%M%S)

all: export GOPATH=${PWD}/../../../..
all: format
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -i -o bin/trustless github.com/untoldwind/trustless
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -i -o bin/trustless-native github.com/untoldwind/trustless/native

install.local: export GOPATH=${PWD}/../../../..
install.local: all
	@cp bin/trustless ${HOME}/bin
	@sed 's:@@@HOME@@@:'"${HOME}"':g' scripts/trustless.service > ${HOME}/.config/systemd/user/trustless.service
	@systemctl --user daemon-reload

format: export GOPATH=${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ${PKGS}

test: export GOPATH=${PWD}/../../../..
test:
	@echo "--> Running tests"
	@go test -v ${PKGS}

cross: bin.linux64 bin.macos bin.windows64 bin.windows32

bin.linux64: export GOPATH=${PWD}/../../../..
bin.linux64: export GOOS=linux
bin.linux64: export GOARCH=amd64
bin.linux64: export CGO_ENABLED=1
bin.linux64:
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-linux-amd64 github.com/untoldwind/trustless
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -i -o bin/trustless-native-linux-amd64 github.com/untoldwind/trustless/native

bin.macos: export GOPATH=${PWD}/../../../..
bin.macos: export GOOS=darwin
bin.macos: export GOARCH=amd64
bin.macos:
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-darwin-amd64 github.com/untoldwind/trustless
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-native-darwin-amd64 github.com/untoldwind/trustless/native

bin.windows64: export GOPATH=${PWD}/../../../..
bin.windows64: export GOOS=windows
bin.windows64: export GOARCH=amd64
bin.windows64: export CGO_ENABLED=1
bin.windows64:
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-windows-amd64.exe github.com/untoldwind/trustless
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-native-windows-amd64.exe github.com/untoldwind/trustless/native

bin.windows32: export GOPATH=${PWD}/../../../..
bin.windows32: export GOOS=windows
bin.windows32: export GOARCH=386
bin.windows32: export CGO_ENABLED=1
bin.windows32:
	@mkdir -p bin
	@echo "--> Running go build ${VERSION}"
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-windows-x86.exe github.com/untoldwind/trustless
	@go build -ldflags "-w -X github.com/untoldwind/trustless/config.version=${VERSION}" -v -o bin/trustless-native-windows-x86.exe github.com/untoldwind/trustless/native

glide.install:
	@echo "--> glide install"
	@go get github.com/Masterminds/glide
	@go build -v -o bin/glide github.com/Masterminds/glide
	@bin/glide install -v

release: cross
	@echo "--> github-release"
	@go get github.com/c4milo/github-release
	@go build -v -o bin/github-release github.com/c4milo/github-release
	@bin/github-release 21re/trustless ${VERSION} master ${VERSION} 'bin/trustless-*'
