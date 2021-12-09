#!/bin/bash

rm go.mod
rm go.sum

go mod init trade_agent
go mod tidy
