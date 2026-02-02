FROM golang:1.22-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN go build -o main ./api/main.go


FROM alpine:latest

WORKDIR /root/


RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY .env . 

EXPOSE 8080

CMD ["./main"]