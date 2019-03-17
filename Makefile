TAG := $(shell git tag -l --points-at @)

setup:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	go get -u golang.org/x/tools/cmd/cover
	npm install
	dep ensure

build:
	npm run build

test:
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
	cd docs && hugo && cd ..

ci: build test codecov release

.PHONY: test docs
