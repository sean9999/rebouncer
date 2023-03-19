#!/bin/sh

echo input semver
read version

echo GOPROXY=proxy.golang.org go list -m github.com/sean9999/rebouncer@${version}

