export JFROG_CLI_HOME_DIR=${PWD}/.jfrog

build:
	go build -o build/jfrogsetmeup 

test:
	go test ./...

.PHONY: build test