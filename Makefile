test: 
	go test -v ./...

api:
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

test/integration:
	docker compose -f test/integration/docker-compose.yaml up -d --wait
	bash -c "trap '$(MAKE) test/integration/down' EXIT; $(MAKE) test/integration/run"

test/integration/down:
	docker compose -f test/integration/docker-compose.yaml down -v

test/integration/run:
	POSTGRES_PASSWORD=integration \
	USER_SERVICE_SYMMETRIC_KEY=12345678901234567890123456789012 \
	go test -tags integration -v --coverprofile=coverage.out -coverpkg ./... ./test/integration

lint:
	golangci-lint run --fix
