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
	@$(GOLANG) get github.com/go-ldap/ldap/v3
#	@$(GOLANG) get gopkg.in/ldap.v3
	@$(GOLANG) get github.com/go-ldap/ldap/v3
	@$(GOLANG) get golang.org/x/sys/unix
	@$(GOLANG) get gopkg.in/yaml.v2
	@go get module/go-daemon
#	@cd ./src/core/utils/; $(GOLANG) build utils.go; $(GOLANG) install;
#	@cd ./src/core/logs/; $(GOLANG) build logs.go; $(GOLANG) install;
	@cd ./src/module/go-daemon/; $(GOLANG) build .; $(GOLANG) install;
	@cd ./src/module/tacplus/; $(GOLANG) build .; $(GOLANG) install;
	@cd ./src/module/ldap-client/; $(GOLANG) build .; $(GOLANG) install;
	@$(GOLANG) build tacacs.go libs.go tacacs_libs.go ldap_auth.go

.PHONY: help

all: help
