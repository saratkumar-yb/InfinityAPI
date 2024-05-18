FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o infinityapi .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/infinityapi /app/infinityapi
COPY --from=builder /app/config.ini /app/config.ini
COPY --from=builder /app/schema.sql /app/schema.sql

COPY scripts/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
