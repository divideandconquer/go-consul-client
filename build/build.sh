#!/bin/bash
CGO_ENABLED=0 godep go build -a -ldflags "-s" -installsuffix cgo -o bin/app src/main/*.go