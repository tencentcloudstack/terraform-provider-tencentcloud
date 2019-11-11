# terraform-provider-tencentcloud

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x
* [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Usage

### Build from source code

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-tencentcloud`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/terraform-providers/terraform-provider-tencentcloud
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

### Terrafrom it

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
How to trigger running the test cases, please refer the `test.sh` script.
How to write test cases, check the `xxx_test.go` files.

### Avoid ``terrafrom init``

```
export TF_SKIP_PROVIDER_VERIFY=1
```

This will disable the verify steps, so after you update this provider, you won't need to create new resources, but use previously saved state.

### Document

Keep in mind that document changes is also needed when resources, data sources, attributes changed in code.
