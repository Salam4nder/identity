FROM golang:latest as builder

WORKDIR /app

COPY . .
COPY db/migrations /app/db/migrations

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main ./cmd/app/

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations /app/db/migrations
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8081

CMD ["./main"]
