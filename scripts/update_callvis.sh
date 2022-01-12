#!/bin/bash

go-callvis \
    -group pkg,type \
    -minlen 5 \
    -focus trade_agent/pkg/modules/order \
    -skipbrowser \
    -file=./assets/callvis \
    ./cmd || exit 1
