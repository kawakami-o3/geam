all: build

build:
	goimports -w -l .
	go generate ./...
	go build

test: build
	go test

clean:
	go clean



