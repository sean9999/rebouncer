#!/bin/sh

echo input semver
read version

GOPROXY=proxy.golang.org go list -m github.com/sean9999/rebouncer@${version}

