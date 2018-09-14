---
layout: "tencentcloud"
page_title: "Provider: tencentcloud"
sidebar_current: "docs-tencentcloud-index"
description: |-
  The TencentCloud provider is used to interact with many resources supported by TencentCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# TencentCloud Provider

The TencentCloud provider is used to interact with the
many resources supported by [TencentCloud](https://intl.cloud.tencent.com). The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = "${var.secret_id}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

# Create a web server
resource "tencentcloud_instance" "web" {
  instance_name              = "ruby on rails example"
  availability_zone          = "ap-guangzhou-3"
  image_id                   = "img-xxxxxxxx"
  instance_type              = "S1"
  key_name                   = "${tencentcloud_key_pair.my_ssh_key.id}"
  security_groups            = ["${tencentcloud_security_group.default.id}"]
  internet_max_bandwidth_out = 20
  count                      = 1
}

# Create key pair with your public key
resource "tencentcloud_key_pair" "my_ssh_key" {
  key_name = "from_terraform_public_key"
  public_key = "ssh-rsa AAAAB3NzaSuperLongString foo@bar"
}

# Create security group
// Create Security Group with 2 rules
resource "tencentcloud_security_group" "default" {
  name        = "web accessibility"
  description = "make it accessable for both production and stage ports"
}
resource "tencentcloud_security_group_rule" "web" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,3000,8080"
  policy            = "accept"
}
resource "tencentcloud_security_group_rule" "ssh" {
  security_group_id = "${tencentcloud_security_group.default.id}"
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

Static credentials can be provided by adding an `secret_id` `secret_key` and `region` in-line in the
tencentcloud provider block:

Usage:

```hcl
provider "tencentcloud" {
  secret_id = "${var.secret_id}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```


### Environment variables

You can provide your credentials via `TENCENTCLOUD_SECRET_ID` and `TENCENTCLOUD_SECRET_KEY`,
environment variables, representing your TencentCloud Access Key and Secret Key, respectively.
`TENCENTCLOUD_REGION` is also used, if applicable:

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

The following arguments are supported:

* `secret_id` - (Optional) This is the TencentCloud access key. It must be provided, but
  it can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.

* `secret_key` - (Optional) This is the TencentCloud secret key. It must be provided, but
  it can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.

* `region` - (Required) This is the TencentCloud region. It must be provided, but
  it can also be sourced from the `TENCENTCLOUD_REGION` environment variables.
  The default input value is ap-guangzhou.


## Testing

Credentials must be provided via the `TENCENTCLOUD_SECRET_ID`, and `TENCENTCLOUD_SECRET_KEY` environment variables in order to run acceptance tests.
