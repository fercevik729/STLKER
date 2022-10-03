# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go build -o /stlker-grpc

EXPOSE 9090

ENTRYPOINT ["/stlker-grpc"]

