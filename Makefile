.PHONY: build run test clean

BINARY_NAME=almari

build:
	go build -o bin/$(BINARY_NAME) main.go

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	rm -rf bin/
