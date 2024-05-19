FROM golang:1.19 AS builder

WORKDIR /app

COPY src/ src/
COPY go.mod .

RUN go build -o build/main src/main/main.go

ENTRYPOINT ["/app/build/main"]
