
.PHONY: clean build run

default: build

setup:
	go get github.com/tools/godep

build:
	godep go build -a -o ./build/atnetgo *.go

install:
	godep go install .

build_darwin:
	GOOS=darwin GOARCH=amd64 godep go build -a -o ./build/atnetgo *.go
	zip ./build/atnetgo_darwin64.zip ./build/atnetgo

build_linux:
	GOOS=linux GOARCH=amd64 godep go build -a -o ./build/atnetgo *.go
	zip ./build/atnetgo_linux64.zip ./build/atnetgo

build_arm5:
	GOOS=linux GOARM=5 GOARCH=arm godep go build -a -o ./build/atnetgo *.go
	zip ./build/atnetgo_linux_arm5.zip ./build/atnetgo

build_arm7:
	GOOS=linux GOARM=7 GOARCH=arm godep go build -a -o ./build/atnetgo *.go
	zip ./build/atnetgo_linux_arm7.zip ./build/atnetgo

build_win64:
	GOOS=windows GOARCH=amd64 godep go build -a -o ./build/atnetgo.exe *.go
	zip ./build/atnetgo_win64.zip ./build/atnetgo.exe

build_win32:
	GOOS=windows GOARCH=386 godep go build -a -o ./build/atnetgo.exe *.go
	zip ./build/atnetgo_win32.zip ./build/atnetgo.exe

all: build_darwin build_linux build_arm5 build_arm7 build_win64 build_win32
	rm ./build/atnetgo
	rm ./build/atnetgo.exe

run: build
	./build/atnetgo

clean:
	- rm -r build