#!/usr/bin/env bash

# This is a screen script to build the go application
# It takes the os and arch as arguments
# The default assumptions are linux amd64

# Set the default values
OS=${1:-linux}
ARCH=${2:-amd64}
GOOS=OS \
GOARCH=GOARCH \
CGO_ENABLED=1 \
CC="zig cc -target x86_64-windows" \
CXX="zig c++ -target x86_64-windows" \
# Build the application
GOARCH=$ARCH GOOS=$OS \
go build -trimpath -ldflags='-H=windowsgui -r' -o bin/${OS}/${ARCH}/