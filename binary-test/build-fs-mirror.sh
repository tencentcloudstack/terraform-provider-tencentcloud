#!/bin/bash

# This script construct local filesystem provider mirror
# To make binary work, set the environment variable: TF_CONFIG_FILE=$PWD/.terraform

set -eo pipefail

provider_name=terraform-provider-tencentcloud
# using `localtest` to make sure the provider fetch from local
namespace=localtest
os=$(uname -s | tr "[:upper:]" "[:lower:]")
arch=$(uname -m)
os_arch="${os}_${arch}"
version=$1
#if [ -z $version ]; then
#  echo "missing argument: tag_version"
#  echo "example usage: ./scripts/build-artifact 1.0.0"
#  exit -1
#fi

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
