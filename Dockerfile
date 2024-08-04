FROM golang:latest as builder

WORKDIR /app

COPY . .
COPY /internal/database/migrations /app/db/migrations
COPY ./config.yaml .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main .

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations /app/db/migrations
COPY --from=builder /app/config.yaml .
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 50051 8080

CMD ["./main"]
