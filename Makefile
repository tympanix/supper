TAG := $(shell git tag -l --points-at @)

setup:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/gobuffalo/packr/v2/packr2
	npm install
	dep ensure

build:
	npm run build

clean:
	rm -rf web/build
	rm -rf docs/public
	rm -rf dist
	packr2 clean

mest:
	go test -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...

prerelease:
	packr2

release: prerelease
ifdef TAG
	git status
	curl -sL http://git.io/goreleaser | bash
else
	@echo "Skip publish..."
endif

codecov:
	curl -sL https://codecov.io/bash | bash

docs:
	rm -rf docs/public
	git worktree prune
	git worktree add docs/public gh-pages
	git submodule update --recursive --remote
	cd docs && hugo && cd ..

ci: build test codecov release

.PHONY: test docs clean
