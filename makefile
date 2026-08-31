.PHONY: all server tidy test clean docker-up docker-down

all: server

server:
	go run cmd/api/main.go

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v --remove-orphans

tidy:
	go mod tidy
test:
	go test ./..
clean:
	go clean -cache -testcache -modcache