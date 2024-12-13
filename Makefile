build:
	go build -o bin/main cmd/server/main.go

run:
	go run cmd/server/main.go

migrate:
	go run cmd/migration/main.go

test:
	go test ./internal/handlers

lint:
	golangci-lint run