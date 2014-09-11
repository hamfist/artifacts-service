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
all: clean test save

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
	$(GO) test -covermode=count -coverprofile=$@.tmp $(GOBUILD_LDFLAGS) $(PACKAGE)
	echo 'mode: count' > $@
	grep -h -v 'mode: count' $@.tmp >> $@
	rm -f $@.tmp
	grep -h -v 'mode: count' $^ >> $@
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
	$(RM) -vf $(shell godep path)/bin/artifacts-service $(GOPATH)/bin/artifacts-service
	$(RM) -vf coverage.html *coverage.coverprofile
	$(GO) clean $(PACKAGE) $(SUBPACKAGES) || true
	if [ -d $(shell godep path)/pkg ] ; then \
		$(FIND) $(shell godep path)/pkg -wholename \
			'*hamfist/artifacts-service*' | xargs rm -rfv || true ; \
	fi ; \
	if [ -d $(GOPATH)/pkg ] ; then \
		$(FIND) $(GOPATH)/pkg -wholename \
			'*hamfist/artifacts-service*' | xargs rm -rfv || true ; \
	fi

.PHONY: save
save:
	$(GODEP) save

.PHONY: fmtpolice
fmtpolice:
	set -e; $(foreach f,$(shell git ls-files '*.go' | grep -v Godeps),gofmt $(f) | diff -u $(f) - ;)

.PHONY: lintall
lintall:
	set -e; golint $(PACKAGE) ; $(foreach pkg,$(SUBPACKAGES),golint $(pkg) ;)
