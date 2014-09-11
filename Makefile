PACKAGE := github.com/hamfist/artifacts-service
SUBPACKAGES := \
	$(PACKAGE)/artifact \
	$(PACKAGE)/auth \
	$(PACKAGE)/metadata \
	$(PACKAGE)/server \
	$(PACKAGE)/store

VERSION_VAR := main.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)

REV_VAR := main.RevisionString
REPO_REV := $(shell git rev-parse --sq HEAD)

FIND ?= find
GO ?= godep go
GODEP ?= godep
GOPATH := $(shell echo $${GOPATH%%:*})
GOBUILD_LDFLAGS := -ldflags "-X $(VERSION_VAR) $(REPO_VERSION) -X $(REV_VAR) $(REPO_REV)"
GOBUILD_FLAGS ?= -x

DATABASE_URL ?= 'postgres://artifacts:dogs@localhost:5432/artifacts?sslmode=disable'
PORT ?= 9839

export DATABASE_URL
export PORT

COVERPROFILES := \
	artifact-coverage.coverprofile \
	auth-coverage.coverprofile \
	metadata-coverage.coverprofile \
	server-coverage.coverprofile \
	store-coverage.coverprofile

.PHONY: all
all: clean test lintall

.PHONY: test
test: build fmtpolice test-deps test-race coverage.html

.PHONY: test-deps
test-deps:
	$(GO) test -i $(GOBUILD_LDFLAGS) $(PACKAGE) $(SUBPACKAGES)

.PHONY: test-race
test-race:
	$(GO) test -race $(GOBUILD_LDFLAGS) $(PACKAGE) $(SUBPACKAGES)

coverage.html: coverage.coverprofile
	$(GO) tool cover -html=$^ -o $@

coverage.coverprofile: $(COVERPROFILES)
	./bin/fold-coverprofiles $^ > $@
	$(GO) tool cover -func=$@

artifact-coverage.coverprofile:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/artifact

auth-coverage.coverprofile:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/auth

metadata-coverage.coverprofile:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/metadata

server-coverage.coverprofile:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/server

store-coverage.coverprofile:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/store

.PHONY: build
build:
	$(GO) install $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(PACKAGE)

.PHONY: deps
deps:
	$(GO) get $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(PACKAGE)

.PHONY: clean
clean:
	./bin/clean

.PHONY: save
save:
	$(GODEP) save

.PHONY: fmtpolice
fmtpolice:
	./bin/fmtpolice

.PHONY: lintall
lintall:
	./bin/lintall
