#
# simple makefile to build and release k2env
#

PWD                       := $(shell pwd)
PREFIX                    ?= $(GOPATH)
BINDIR                    ?= $(PREFIX)/bin
GO                        := GO111MODULE=on go
# GOOS                      ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f1)
# GOARCH                    ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f2)


build: build-linux build-macos build-windows

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o azapim.linux cmd/azapim/azapim.go

build-macos:
	GOOS=darwin GOARCH=amd64 $(GO) build -o azapim.macos cmd/azapim/azapim.go

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o azapim.windows cmd/azapim/azapim.go
