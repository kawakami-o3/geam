all: build

build:
	#goimports -w -l .
	go generate ./...
	#go build
	env GO111MODULE=on go build

test: build
	go test

clean:
	#go clean
	env GO111MODULE=on go clean -cache

deps:
	go get github.com/kawakami-o3/go-genopc

