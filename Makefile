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
	curl -sL http://git.io/goreleaser | bash
else
	@echo "Skip publish..."
endif

codecov:
	bash <(curl -s https://codecov.io/bash)

ci: setup build test codecov release
