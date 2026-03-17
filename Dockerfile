# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download || go mod download # Retry once for stability

COPY . .

RUN go build -o shipment-server cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/shipment-server .

EXPOSE 50051

CMD ["./shipment-server"]
