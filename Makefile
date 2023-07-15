build:
	go build cmd/envserver.go

run: build
	./envserver

test: build
	go test ./...

clean:
	rm -rf ./envserver