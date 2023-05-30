#!/bin/bash

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <OS> <ARCH>"
  echo "Example: $0 windows amd64"
  exit 1
fi

GOOS=$1
GOARCH=$2
flags=""
output="./build/g45w_${GOOS}_${GOARCH}"

if [ $GOOS = "windows" ]; then
  output+=".exe"
  flags="-H windowsgui"
fi

GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$flags" -o "$output" .
