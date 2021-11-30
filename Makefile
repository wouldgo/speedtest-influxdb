SHELL := /bin/bash
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
OUT := $(shell pwd)/_out

zsh-autocomplete: deps
	go run main.go --zsh-autocomplete

clean-compile: clean deps compile

compile:
	CGO_ENABLED=0 \
	GOOS=$(OS) \
	GOARCH=$(ARCH) \
	go build \
	-ldflags '-w -extldflags "-static"' \
	-a -o $(OUT)/speedtest *.go

deps:
	go mod download

clean:
	rm -Rf $(OUT)
	mkdir -p $(OUT)
	touch $(OUT)/.keep
