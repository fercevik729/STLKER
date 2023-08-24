.PHONY: go

DIR := ./protos

go:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./protos/watcher.proto

js:
	protoc --js_out=import_style=commonjs:. ./protos/watcher.proto
	protoc --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. ./protos/watcher.proto
