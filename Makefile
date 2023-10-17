config ?= ""

build:
	go build cmd/envserver.go

run: build
	./envserver -config $(config)

test: build
	go clean -testcache && go test ./...

clean:
	rm -rf ./envserver

install:
	go mod download