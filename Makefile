
.PHONY: clean build run

default: build

default: build_darwin

setup:
	go get github.com/tools/godep
	
build_darwin: 
	GOOS=darwin GOARCH=amd64 godep go build -a -o ./build/atnetgo *.go

build_linux: 
	GOOS=linux GOARCH=amd64 godep go build -a -o ./build/atnetgo *.go

run: build
	./build/atnetgo