NAME     ?= ddgo
VERSION  ?= $(shell git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/[-+].*$$//')
REVISION ?= $(shell git rev-parse --short HEAD)
BUILD_OS_TARGETS   = "linux darwin windows"
BUILD_ARCH_TARGETS = "amd64"

BUILD_LDFLAGS = "-s -w"

all: clean test build

test: lint
	go test -v -short $(TESTFLAGS) ./...

deps:
	go get -d -v -t ./...
	go get golang.org/x/tools/cmd/goimports
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/golang/lint/golint
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls
	go get github.com/laher/goxc
	
lint: deps
	go tool vet -all -printfuncs=Criticalf,Infof,Warningf,Debugf,Tracef .
	goimports -l .
	megacheck .
	_tools/go-linter $(BUILD_OS_TARGETS)

xbuild: deps
	goxc -build-ldflags=$(BUILD_LDFLAGS) \
		-os=$(BUILD_OS_TARGETS) -arch=$(BUILD_ARCH_TARGETS) -d . -n $(NAME)