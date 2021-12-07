#!/bin/bash

rm go.mod
rm go.sum

go mod init gitlab.tocraw.com/root/toc_trader
go mod tidy
