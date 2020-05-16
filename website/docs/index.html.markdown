---
layout: "tencentcloud"
page_title: "Provider: tencentcloud"
sidebar_current: "docs-tencentcloud-index"
description: |-
  The TencentCloud provider is used to interact with many resources supported by TencentCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# TencentCloud Provider

The TencentCloud provider is used to interact with many resources supported by [TencentCloud](https://intl.cloud.tencent.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** Terraform 0.12.x support began with provider version 1.9.0 (June 18, 2019), 

## Example Usage

```hcl
# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
}

# Create a web server
resource "tencentcloud_instance" "web" {
  instance_name              = "web example"
  availability_zone          = "ap-guangzhou-3"
  image_id                   = "img-9qabwvbn"
  instance_type              = "S1.SMALL1"
  system_disk_type           = "CLOUD_PREMIUM"
  key_name                   = tencentcloud_key_pair.my_ssh_key.id
  security_groups            = [tencentcloud_security_group.default.id]
  internet_max_bandwidth_out = 20
  count                      = 1
}

# Create key pair with your public key
resource "tencentcloud_key_pair" "my_ssh_key" {
  key_name   = "from_terraform_public_key"
  public_key = "ssh-rsa AAAAB3NzaSuperLongString foo@bar"
}

# Create security group with 2 rules
resource "tencentcloud_security_group" "default" {
  name        = "web accessibility"
  description = "make it accessible for both production and stage ports"
}

resource "tencentcloud_security_group_rule" "web" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,3000,8080"
  policy            = "accept"
}

resource "tencentcloud_security_group_rule" "ssh" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = "202.119.230.10/32"
  ip_protocol       = "tcp"
  port_range        = "22"
  policy            = "accept"
}
```

## Authentication

The TencentCloud provider offers a flexible means of providing credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials ###

Static credentials can be provided by adding an `secret_id` `secret_key` and `region` in-line in the tencentcloud provider block:

Usage:

```hcl
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
}
```

### Environment variables

You can provide your credentials via `TENCENTCLOUD_SECRET_ID` and `TENCENTCLOUD_SECRET_KEY`, environment variables,
representing your TencentCloud Access Key and Secret Key, respectively. `TENCENTCLOUD_REGION` is also used, if applicable:

```hcl
provider "tencentcloud" {}
```

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="your_fancy_accesskey"
$ export TENCENTCLOUD_SECRET_KEY="your_fancy_secretkey"
$ export TENCENTCLOUD_REGION="ap-guangzhou"
$ terraform plan
```

## Argument Reference

In addition to generic provider arguments (e.g. alias and version), the following arguments are supported in the TencentCloud provider block:

* `secret_id` - This is the TencentCloud access key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.
* `secret_key` - This is the TencentCloud secret key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.
* `security_token` - TencentCloud Security Token of temporary access credentials. It can be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).
* `region` - This is the TencentCloud region. It must be provided, but it can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is ap-guangzhou.

## Testing

Credentials must be provided via the `TENCENTCLOUD_SECRET_ID`, and `TENCENTCLOUD_SECRET_KEY` environment variables in order to run acceptance tests.
