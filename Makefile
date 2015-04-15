GO ?= go
GOPATH := $(CURDIR)/_vendor:$(GOPATH)

default: build_client build_server

build_client:
	./vendor.sh
	cd $(CURDIR)/client && $(GO) build

build_server:
	./vendor.sh
	cd $(CURDIR)/server && $(GO) build