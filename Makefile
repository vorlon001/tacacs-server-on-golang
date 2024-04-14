-include .env
#VERSION := $(shell git describe --tags)
#BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
GOLANG := "/usr/local/go/bin/go"
PWDSRC := $(shell pwd)

## help: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: BUILD
build:
	@echo "BUILD"
	@export GOPATH=$GOPATH:$(PWDSRC)
	@echo "VARS MAKE: $(GOLANG) $(PWDSRC) $(GOPATH)"
	@$(GOLANG) mod tidy -e
	@CGO_ENABLED=0 $(GOLANG) build -ldflags "-w -s -X 'main.Version=1.0.3'" -o tacacsserver tacacs.go libs.go tacacs_libs.go ldap_auth.go tacacs-server.go

.PHONY: help

all: help
