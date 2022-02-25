#!/bin/bash

git clone git@gitlab.tocraw.com:root/trade_agent_protobuf.git

protoc -I=. --go_out=. ./trade_agent_protobuf/src/*.proto

rm -rf trade_agent_protobuf
