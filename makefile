export JFROG_CLI_HOME_DIR=$(CURDIR)/.jfrog

build:
	go build -o build/jfrogsetmeup 

test:
	go test ./... -p=1 -count=1

coverage:
	go test ./... -coverprofile=./build/coverage.out

open-coverage: coverage
	go tool cover -html=./build/coverage.out

calc-coverage: coverage
	go tool cover -func=./build/coverage.out

coverage-badge:
	gopherbadger -covercmd "make calc-coverage"
	mv ./coverage_badge.png ./assets/images/coverage.png

.PHONY: build test