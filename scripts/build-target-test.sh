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

provider_name=terraform-provider-tencentcloud
# using `localtest` to make sure the provider fetch from local
namespace=localtest
os=$(uname -s | tr "[:upper:]" "[:lower:]")
arch=$(uname -m)
os_arch="${os}_${arch}"

echo "building $provider_name $version filesystem_mirror"

plugin_dir=".terraform/registry.terraform.io/$namespace/tencentcloud"

rm -rf .terraform

mkdir -p $plugin_dir

go build

zip "$plugin_dir/${provider_name}_${version}_${os_arch}.zip" $provider_name

tee -a $plugin_dir/index.json <<EOF
{
  "versions": {
    "$version": {}
  }
}
EOF

tee -a "$plugin_dir/${version}.json" <<EOF
{
  "archives": {
    "${os_arch}": {
      "url": "${provider_name}_${version}_${os_arch}.zip"
    }
  }
}
EOF

tee -a .terraform/prerelease.tfrc <<EOF
provider_installation {
	 filesystem_mirror {
         path    = "$PWD/.terraform"
         include = ["registry.terraform.io/$namespace/tencentcloud"]
     }
}
EOF

echo $PWD

go test -v -run TestTerraformBasicExample ./tencentcloud/binary-test -version=$version -source=localtest/tencentcloud