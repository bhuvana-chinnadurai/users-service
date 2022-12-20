# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .

ARG SERVER_DIR=api/server

RUN go mod download

RUN	go build -o ./api-server ./${SERVER_DIR}/server.go

EXPOSE 8080

CMD ["./api-server"]