.PHONY: all server tidy test clean

all: server

server:
	go run cmd/api/main.go
tidy:
	go mod tidy
test:
	go test ./..
clean:
	go clean -cache -testcache -modcache