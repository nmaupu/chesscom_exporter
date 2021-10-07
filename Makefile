BIN=bin
BIN_NAME=chesscom-exporter

PKG_NAME = github.com/nmaupu/chesscom_exporter
TAG_NAME ?= $(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)
LDFLAGS = -ldflags="-X 'main.AppVersion=$(TAG_NAME)' -X 'main.BuildDate=$(shell date)'"

all: $(BIN_NAME)

$(BIN_NAME): $(BIN)
	go build -o $(BIN)/$(BIN_NAME) $(LDFLAGS)

.PHONY: release
release:
	GOOS=windows GOARCH=amd64 go build -o $(BIN)/$(BIN_NAME)-$(TAG_NAME)-windows_x64.exe $(LDFLAGS)
	GOOS=windows GOARCH=386   go build -o $(BIN)/$(BIN_NAME)-$(TAG_NAME)-windows_x86.exe $(LDFLAGS)
	GOOS=darwin GOARCH=amd64  go build -o $(BIN)/$(BIN_NAME)-$(TAG_NAME)-darwin_x64 $(LDFLAGS)

$(BIN):
	mkdir -p $(BIN)

.PHONY: clean
clean:
	rm -rf $(BIN) ./$(BIN_NAME)
