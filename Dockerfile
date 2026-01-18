
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

RUN apk add --no-cache curl
RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

# Открываем порт 8080 наружу
EXPOSE 8080

CMD ["./main"]