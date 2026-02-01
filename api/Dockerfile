# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN go build -o main ./api/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

# Install certificates for AWS SSL connection
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY .env . 

EXPOSE 8080

CMD ["./main"]