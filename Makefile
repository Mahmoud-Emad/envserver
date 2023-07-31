build:
	go build cmd/envserver.go

run: build
	./envserver

test: build
	go clean -testcache && go test ./...

clean:
	rm -rf ./envserver

install:
	go mod download