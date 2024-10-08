.PHONY: test test-cover test-db test-db/down test-db/run run api up down proto lint nancy

test: 
	go test -count=1 ./... && $(MAKE) test-db

test-cover:
	go tool cover -html=coverage.out

test-db:
	docker compose -f internal/database/docker-compose.yaml up -d --wait
	bash -c "trap '$(MAKE) test-db/down' EXIT; $(MAKE) test-db/run"

test-db/down:
	docker compose -f internal/database/docker-compose.yaml down -v

test-db/run:
	go test -count=1 -tags testdb --coverprofile=coverage.out -coverpkg ./... ./internal/database/...

api:
	docker build -t identity .

up:
	make api && docker compose up -d

down:
	docker compose down -v

proto:
	rm -rf proto/gen/*.go
	protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
     proto/*.proto

lint:
	docker run -t --rm -v $(shell pwd):/app -v ~/.cache/golangci-lint/v1.57.2:/root/.cache -w /app golangci/golangci-lint:v1.57.2 golangci-lint run -v

nancy:
	go list -json -deps ./... | docker run --rm -i sonatypecommunity/nancy:latest sleuth
