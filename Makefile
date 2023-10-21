.PHONY: test test-cover test-db test-db/down test-db/run run docker up down logs logs-user logs-db evans proto lint nancy

test: 
	go test -v ./...

test-cover:
	go tool cover -html=coverage.out

test-db:
	docker compose -f internal/db/docker-compose.yaml up -d --wait
	bash -c "trap '$(MAKE) test-db/down' EXIT; $(MAKE) test-db/run"

test-db/down:
	docker compose -f internal/db/docker-compose.yaml down -v

test-db/run:
	go test -tags testdb -v --coverprofile=coverage.out -coverpkg ./... ./internal/db

run:
	go run ./cmd/app/main.go

docker:
	docker build -t user .

up:
	docker-compose up -d

down:
	docker-compose down -v

logs:
	docker-compose logs -f

logs-user:
	docker-compose logs -f api

logs-db:
	docker-compose logs -f postgres

evans:
	evans -r
	
proto:
	rm -rf internal/proto/gen/*.go
	protoc --proto_path=internal/proto --go_out=internal/proto/gen --go_opt=paths=source_relative \
    --go-grpc_out=internal/proto/gen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/proto/gen --grpc-gateway_opt=paths=source_relative \
     internal/proto/*.proto

lint:
	golangci-lint run

nancy:
	go list -json -deps ./... | docker run --rm -i sonatypecommunity/nancy:latest sleuth
