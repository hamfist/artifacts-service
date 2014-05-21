PACKAGE := github.com/meatballhat/artifacts-service
SUBPACKAGES := $(PACKAGE)/server $(PACKAGE)/store

VERSION_VAR := main.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)

REV_VAR := main.RevisionString
REPO_REV := $(shell git rev-parse --sq HEAD)

GO ?= go
GODEP ?= godep
GOBUILD_LDFLAGS := -ldflags "-X $(VERSION_VAR) $(REPO_VERSION) -X $(REV_VAR) $(REPO_REV)"
GOBUILD_FLAGS ?=

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

coverage.html: coverage.out
	$(GO) tool cover -html=$^ -o $@

coverage.out: server-coverage.out store-coverage.out artifact-coverage.out
	$(GO) test -covermode=count -coverprofile=$@.tmp $(GOBUILD_LDFLAGS) $(PACKAGE)
	echo 'mode: count' > $@
	grep -h -v 'mode: count' $@.tmp >> $@
	rm -f $@.tmp
	grep -h -v 'mode: count' $^ >> $@
	$(GO) tool cover -func=$@

server-coverage.out:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/server

store-coverage.out:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/store

artifact-coverage.out:
	$(GO) test -covermode=count -coverprofile=$@ $(GOBUILD_LDFLAGS) $(PACKAGE)/artifact

.PHONY: build
build: deps
	$(GO) install $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(PACKAGE)

.PHONY: deps
deps:
	$(GO) get $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(PACKAGE)
	$(GODEP) restore

.PHONY: clean
clean:
	rm -vf $${GOPATH%%:*}/bin/artifacts-service
	rm -vf coverage.html *coverage.out
	$(GO) clean $(PACKAGE) $(SUBPACKAGES) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -name '*artifacts-service*' | xargs rm -rfv || true ; \
	fi

.PHONY: save
save:
	$(GODEP) save -copy=false

.PHONY: fmtpolice
fmtpolice:
	set -e; for f in $(shell git ls-files '*.go'); do gofmt $$f | diff -u $$f - ; done
