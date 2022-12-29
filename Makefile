APPNAME = pi-kit
BIN = $(GOPATH)/bin
GOCMD = go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)/pi-kit
BINARY_WINDOWS = $(BIN)/pi-kit.exe
BUILD_TIME = $(shell date '+%Y-%m-%d %H:%M:%S')
PID = .pid
HUB_ADDR = hub.nsini.com
DOCKER_USER =
DOCKER_PWD =
VERSION = $(shell git describe --tags --always --dirty)
GO_LDFLAGS = -ldflags="-X 'github.com/pi-kit/pi-kit/cmd/service.version=$(VERSION)' -X 'github.com/pi-kit/pi-kit/cmd/service.buildDate=$(BUILD_TIME)' -X 'github.com/pi-kit/pi-kit/cmd/service.gitCommit=$(shell git rev-parse --short HEAD)' -X 'github.com/pi-kit/pi-kit/cmd/service.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)'"
NAMESPACE = pi-kit
PWD = $(shell pwd)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

start:
	$(BIN)/$(APPNAME) -v
	$(BIN)/$(APPNAME) start -p :8080 & echo $$! > $(PID)

restart:
	@echo restart the app...
	@kill `cat $(PID)` || true
	$(BIN)/$(APPNAME) start -p :8080 & echo $$! > $(PID)

stop:
	@kill `cat $(PID)` || true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_UNIX) $(GO_LDFLAGS) ./cmd/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -v -o $(BINARY_WINDOWS) $(GO_LDFLAGS) ./cmd/main.go

login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PWD) $(HUB_ADDR)

build-docker:
	docker build --rm -t $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(VERSION) .

push-docker:
	docker push $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(VERSION)

build:
	CGO_ENABLED=0 $(GOBUILD) -v -o $(BINARY_UNIX) $(GO_LDFLAGS) ./cmd/main.go

init:
	GO111MODULE=on $(GORUN) ./cmd/main.go generate table all

run:
	GO111MODULE=on $(GORUN) ./cmd/main.go start