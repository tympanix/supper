TAG := $(shell git tag -l --points-at @)

setup:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/rakyll/statik
	npm install
	dep ensure

build:
	npm run build
	go generate

clean:
	rm -rf web/build
	rm -rf docs/public
	rm -rf dist

mest:
	go test -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...

release:
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
	git submodule update --init --recursive
	cd docs && hugo && cd ..

dist: clean build
	goreleaser release --skip-publish --skip-validate

ci: build test codecov release

.PHONY: test docs clean
