#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -o ./bin/kblswitch.exe ./cmd/main.go
