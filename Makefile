.PHONY: build build-frontend run test clean

BINARY_NAME=almari

build-frontend:
	cd frontend && npm ci && npm run build

build: build-frontend
	go build -o bin/$(BINARY_NAME) main.go

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

test-frontend:
	cd frontend && npm run test

clean:
	rm -rf bin/ static/

dev-frontend:
	cd frontend && npm run dev

dev: build-frontend
	./bin/$(BINARY_NAME)