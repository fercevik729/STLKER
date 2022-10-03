# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /stlker-grpc

EXPOSE 9090

ENTRYPOINT ["/stlker-grpc"]

