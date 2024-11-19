# Makefile
.PHONY: build run down clean docker-build docker-run go-run go-build

build:
	docker-compose build

run:
	docker-compose up

down:
	docker-compose down

clean:
	docker system prune -af

go-run:
	go run ./cmd/api

go-build:
	go build ./cmd/api

go-get:
	go get ./...
	
test:
	go test ./...

