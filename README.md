<a href="https://terraform.io">
   <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
</a>

# Terraform Provider For TencentCloud

[![stars](https://img.shields.io/github/stars/tencentcloudstack/terraform-provider-tencentcloud)](https://img.shields.io/github/stars/tencentcloudstack/terraform-provider-tencentcloud)
[![Forks](https://img.shields.io/github/forks/tencentcloudstack/terraform-provider-tencentcloud)](https://img.shields.io/github/forks/tencentcloudstack/terraform-provider-tencentcloud)
[![Go Report Card](https://goreportcard.com/badge/github.com/tencentcloudstack/terraform-provider-tencentcloud)](https://goreportcard.com/report/github.com/tencentcloudstack/terraform-provider-tencentcloud)
[![Releases](https://img.shields.io/github/release/tencentcloudstack/terraform-provider-tencentcloud.svg?style=flat-square)](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/releases)
[![License](https://img.shields.io/github/license/tencentcloudstack/terraform-provider-tencentcloud)](https://img.shields.io/github/license/tencentcloudstack/terraform-provider-tencentcloud)
[![Issues](https://img.shields.io/github/issues/tencentcloudstack/terraform-provider-tencentcloud)](https://img.shields.io/github/issues/tencentcloudstack/terraform-provider-tencentcloud)

<div>
  <p>
    <a href="https://cloud.tencent.com">
        <img src=".github/01_Tcloud_logo_Eng.png" alt="logo" title="Terraform" height="69">
    </a>
    <br>
    <i>Tencent Infrastructure Automation for Terraform.</i>
    <br>
  </p>
</div>

* Tutorials: https://www.terraform.io
* [![Documentation](https://img.shields.io/badge/documentation-blue)](https://registry.terraform.io/providers/tencentcloudstack/tencentcloud/latest/docs)
* [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
* Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)


## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.13.x
* [Go](https://golang.org/doc/install) 1.18.x (to build the provider plugin)

## Usage

### Build from source code

Clone repository to: `$GOPATH/src/github.com/tencentcloudstack/terraform-provider-tencentcloud`

```sh
$ mkdir -p $GOPATH/src/github.com/tencentcloudstack
$ cd $GOPATH/src/github.com/tencentcloudstack
$ git clone https://github.com/tencentcloudstack/terraform-provider-tencentcloud.git
$ cd terraform-provider-tencentcloud
$ go build .
```

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

## Configuration

### Configure credentials

You will need to have a pair of secret id and secret key to access Tencent Cloud resources, configure it in the provider arguments or export it in environment variables. If you don't have it yet, please access [Tencent Cloud Management Console](https://console.cloud.tencent.com/cam/capi) to create one.

```
export TENCENTCLOUD_SECRET_ID=AKID9HH4OpqLJ5f6LPr4iIm5GF2s-EXAMPLE
export TENCENTCLOUD_SECRET_KEY=72pQp14tWKUglrnX5RbaNEtN-EXAMPLE
```

### Configure proxy info (optional)

If you are beind a proxy, for example, in a corporate network, you must set the proxy environment variables correctly. For example:

```
export http_proxy=http://your-proxy-host:your-proxy-port  # This is just an example, use your real proxy settings!
export https_proxy=$http_proxy
export HTTP_PROXY=$http_proxy
export HTTPS_PROXY=$http_proxy
```

## Run demo

You can edit your own terraform configuration files. Learn examples from examples directory.

Now you can try your terraform demo:

```
terraform init
terraform plan
terraform apply
```

If you want to destroy the resource, make sure the instance is already in ``running`` status, otherwise the destroy might fail.

```
terraform destroy
```

## Developer Guide

### DEBUG

You will need to set an environment variable named ``TF_LOG``, for more info please refer to [Terraform official doc](https://www.terraform.io/docs/internals/debugging.html):

```
export TF_LOG=DEBUG
```

In your source file, import the standard package ``log`` and print the message such as:

```
log.Println("[DEBUG] the message and some import values: %v", importantValues)
```

### Test

The quicker way for development and debug is writing test cases.

Config environment variables:
```
export TF_ACC=true
```

Config your appid for COS bucket testing
```
export TENCENTCLOUD_APPID=1234567890
```

This example show how to test single test function:
```
cd tencentcloud
go test -i; go test -test.run TestAccTencentCloudNatGateway_basic -v
```

To write test cases, check the `xxx_test.go` files for more reference.

### Avoid ``terraform init``

```
export TF_SKIP_PROVIDER_VERIFY=1
```

This will disable the verify steps, so after you update this provider, you won't need to create new resources, but use previously saved state.

### License

Terraform-Provider-TencentCloud is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.