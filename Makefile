# ref. http://postd.cc/auto-documented-makefile/
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

COMMIT := $(shell git describe --always)
VERSION := $(shell cat version.go | perl -ne 'print "v$$1" if /Version = "(.+?)"/')

.PHONY: install
install: ## install dependencies
	go get -u github.com/rhysd/go-github-selfupdate
	go get -u github.com/tcnksm/ghr
	go get -u github.com/mitchellh/gox

.PHONY: release
release: ## release binaries at GitHub (NOTE: update verion.go & the tag before this)
	gox -os 'darwin linux windows' -arch '386 amd64' -ldflags '-X main.GitCommit=$(COMMIT)' -output 'pkg/$(VERSION)/{{.Dir}}_{{.OS}}_{{.Arch}}'
	ghr -u delphinus $(VERSION) pkg/$(VERSION)
