#!/bin/bash

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <OS> <ARCH>"
  echo "Example: $0 windows amd64"
  exit 1
fi

GOOS=$1
GOARCH=$2
OUTPUT="./build/g45w_${GOOS}_${GOARCH}"

source ./build_vars.sh

if [ $GOOS = "windows" ]; then
  OUTPUT+=".exe"
  FLAGS="$FLAGS -H windowsgui"
fi

GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$FLAGS" -o "$OUTPUT" .
