.PHONY: test lint buld clean fmt

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	go fmt ./...

build:
	mkdir -p bin
	go build -o bin/tollgate ./cmd/tollgate
	go build -o bin/shim-git ./cmd/shim-git
	go build -o bin/shim-gh ./cmd/shim-gh

clean:
	rm -rf bin/ dist/ coverage.out coverage.html
