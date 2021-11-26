export JFROG_CLI_HOME_DIR=${PWD}/.jfrog

build:
	go build ./...

test:
	go test ./...
