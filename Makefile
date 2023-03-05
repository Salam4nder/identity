test: 
	go test ./...

run:
	go run ./cmd/app/main.go

up:
	docker-compose up -d

down:
	docker-compose down
