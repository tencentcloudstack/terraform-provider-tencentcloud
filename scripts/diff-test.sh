#!/bin/sh

# You can run `make diff-test` to execute

echo "==> Checking test cases which module modified"

cd ./test-diff

pattern=$(go run main.go | tr -d '\r')

echo "=========== diffs ==========="
git diff --name-only HEAD $(git describe --tags --abbrev=0)
echo "======= end of diffs ========"

if [[ $pattern != "" ]];then
  cd ../tencentcloud

  echo go test -v -run $pattern

  go test -v -run $pattern
else
  echo "No test case match"
fi

