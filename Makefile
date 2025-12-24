# Makefile for PetsProject
export $(shell cat .env | xargs)

APP_NAME = petfeed

.PHONY: all build run test clean docker-up docker-down migrate

all: build

build:
	go build -o $(APP_NAME) ./cmd/petfeed

run:
	./$(APP_NAME)

test:
	go test -v ./...

clean:
	rm -f $(APP_NAME)

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	migrate -database $(DATABASE_URL) -path migrations up
