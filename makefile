export JFROG_CLI_HOME_DIR=$(CURDIR)/.jfrog

build:
	go build -o build/jfrogsetmeup 

test:
	go test ./...

coverage:
	go test ./... -coverprofile=./build/coverage.out

open-coverage: coverage
	go tool cover -html=./build/coverage.out

.PHONY: build test