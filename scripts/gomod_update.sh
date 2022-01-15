#!/bin/bash

rm go.mod
rm go.sum

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/ofabry/go-callvis@latest

go mod init trade_agent
go mod tidy
