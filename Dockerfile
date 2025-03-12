# Этап сборки
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/migrations /root/migrations
COPY --from=builder /app/.env .
COPY --from=builder /app/docs/swagger.json /root/docs/swagger.json

EXPOSE 8080
CMD ["./server"]