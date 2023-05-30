#!/bin/bash

platforms=("windows/amd64" "linux/amd64" "darwin/amd64")

for platform in "${platforms[@]}"
do
  GOOS=$(echo $platform | cut -d "/" -f 1)
  GOARCH=$(echo $platform | cut -d "/" -f 2)
  echo $GOOS

  sh ./build.sh $GOOS $GOARCH
done