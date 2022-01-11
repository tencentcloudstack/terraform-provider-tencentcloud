#!/bin/sh

cd ./test-diff

pattern=$(go run main.go | tr -d '\r')

if [[$pattern != ""]]
  then
    cd ../tencentcloud

    echo go test -v -run $pattern

    go test -v -run $pattern
  else
    echo diff: $(git diff --name-only HEAD $(git describe --tags --abbrev=0))
    echo "no test case match"
fi

