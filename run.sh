#!/bin/sh

GIT_ROOT="$(git rev-parse --show-toplevel)"

go run ./cmd/main.go -dir $GIT_ROOT/tmp
