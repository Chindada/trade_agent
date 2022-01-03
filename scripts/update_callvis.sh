#!/bin/bash

go-callvis \
    -focus main \
    -skipbrowser \
    -file=./assets/callvis \
    ./cmd || exit 1
