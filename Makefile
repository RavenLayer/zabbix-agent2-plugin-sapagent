.ONESHELL:

PACKAGE = zabbix-agent2-plugin-sapagent
BIN_PATH = /usr/sbin/zabbix-agent2-plugin
CONF_PATH = /etc/zabbix/zabbix_agent2.d/plugins.d
TOPDIR := $(CURDIR)

ifeq ($(ARCH), x86)
	GOARCH := 386
else ifeq ($(ARCH), AMD64)
	GOARCH := amd64
else ifeq ($(ARCH), ARM)
	GOARCH := arm
else ifeq ($(ARCH), ARM64)
	GOARCH := arm64
endif

ifndef GOOS
GOOS := $(shell go env GOOS)
endif

ifndef GOARCH
GOARCH := $(shell go env GOARCH)
endif

CGO_ENABLED ?= 0
BUILDMODE   ?= -buildmode=pie
LDFLAGS     ?= -s -w -buildid= -extldflags=-static-pie
GOFLAGS     ?= -trimpath

.PHONY: build install clean lint test format

build:
	go mod tidy
	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(GOOS)" GOARCH="$(GOARCH)" \
		go build $(BUILDMODE) -ldflags="$(LDFLAGS)" $(GOFLAGS) -o $(TOPDIR)/$(PACKAGE)

install: build
	install -m 0755 -d $(BIN_PATH)
	install -m 0755 $(PACKAGE) $(BIN_PATH)/$(PACKAGE)
	install -m 0755 -d $(CONF_PATH)
	install -m 0644 sapagent.conf $(CONF_PATH)/sapagent.conf

clean:
	rm -f $(TOPDIR)/$(PACKAGE)


lint:
	go mod tidy
	golangci-lint run $(TOPDIR)/...

test:
	go mod tidy
	go test -v -tags tests $(TOPDIR)/...

format:
	go fmt $(TOPDIR)/...

