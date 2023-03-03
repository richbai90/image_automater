#!/usr/bin/env bash

# This is a screen script to build the go application
# It takes the os and arch as arguments
# The default assumptions are linux amd64

error() {
    echo "Error: $1"
    exit 1
}

# Set the default values
OS=${1:-linux}
ARCH=${2:-*}
# Check if go is installed
[ $(("$(go env >/dev/null; echo $?)")) -gt 0 ] && error "go is not installed"
# Check if xgo is installed and install it if not
if [ $(("$(xgo -h >/dev/null 2>&1; echo $?)")) -gt 0 ]; then
    echo "xgo not found. Would you like to install it? (y/n)" &&
        read -r install
    if [ "$install" == "y" ]; then
        go install src.techknowlogick.com/xgo@latest
    else
        error "xgo is required for cross platform support"
    fi

fi
# Build the application

mkdir -p "bin/$OS/$ARCH"
xgo --targets="$OS/$ARCH" -out "bin/$OS/" github.com/richbai90/image_automater
rm -rf "bin/$OS/\*"