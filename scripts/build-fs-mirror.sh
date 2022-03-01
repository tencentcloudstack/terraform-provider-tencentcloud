#!/bin/bash

# This script construct local filesystem provider mirror
# To make binary work, set the environment variable: TF_CONFIG_FILE=$PWD/.terraform

provider_name=terraform-provider-tencentcloud
namespace=tencentcloudstack
os=$(uname -s | tr "[:upper:]" "[:lower:]")
arch=$(uname -m)
os_arch="${os}_${arch}"
#tag_name=$(git describe --tags --abbrev=0)
arg=$1
tag_name=${arg:1}

if [ -z $tag_name ]; then
  echo "missing argument: tag_version"
  echo "example usage: ./scripts/build-artifact 1.0.0"
  exit -1
fi

plugin_dir=".terraform/registry.terraform.io/$namespace/tencentcloud"

rm -rf .terraform

mkdir -p $plugin_dir

go build
zip "$plugin_dir/${provider_name}_${tag_name}_${os_arch}.zip" $provider_name

tee -a $plugin_dir/index.json <<EOF
{
  "versions": {
    "$tag_name": {}
  }
}
EOF

tee -a "$plugin_dir/${tag_name}.json" <<EOF
{
  "archives": {
    "${os_arch}": {
      "url": "${provider_name}_${tag_name}_${os_arch}.zip"
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
