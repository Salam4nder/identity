test: 
	go test ./...

run:
	go run ./cmd/app/main.go

up:
	docker-compose up -d

down:
	docker-compose down -v

logs:
	docker-compose logs -f

evans:
	evans -r
	
proto:
	rm -rf internal/proto/gen/*.go
	protoc --proto_path=internal/proto --go_out=internal/proto/gen --go_opt=paths=source_relative \
    --go-grpc_out=internal/proto/gen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/proto/gen --grpc-gateway_opt=paths=source_relative \
     internal/proto/*.proto
