test: 
	go test ./...

run:
	go run ./cmd/app/main.go

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f
	
proto:
	rm -rf internal/proto/pb/*.go
	protoc --proto_path=internal/proto --go_out=internal/proto/pb --go_opt=paths=source_relative \
    --go-grpc_out=internal/proto/pb --go-grpc_opt=paths=source_relative \
    internal/proto/*.proto
