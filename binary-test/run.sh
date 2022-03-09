#!/bin/sh
# this script runs local filesystem provider mirror build and unit test
# run "$ make build-target-test" directly

set -eo pipefail

configFile=$PWD/.terraform/prerelease.tfrc
if [ $TF_CLI_CONFIG_FILE != $configFile  ]; then
  echo "unexpected: EnvVar TF_CLI_CONFIG_FILE"
  echo "expected: $configFile"
  echo "actual:   $TF_CLI_CONFIG_FILE"
  exit -1
fi

latest_tag=$(git describe --tags --abbrev=0)
version=${latest_tag:1}

./binary-test/build-fs-mirror.sh $version

cd ./binary-test

go mod vendor

go test -v -run TestTerraformBasicExample -version=$version -source=localtest/tencentcloud