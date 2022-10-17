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

-> **Note:** Terraform 0.12.x supported began with provider version 1.9.0 (June 18, 2019).

## Example Usage

```hcl
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"
}

# Get availability zones
data "tencentcloud_availability_zones" "default" {
}

# Get availability images
data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

# Get availability instance types
data "tencentcloud_instance_types" "default" {
  cpu_core_count = 1
  memory_size    = 1
}

# Create a web server
resource "tencentcloud_instance" "web" {
  instance_name              = "web server"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 20
  security_groups            = [tencentcloud_security_group.default.id]
  count                      = 1
}

# Create security group
resource "tencentcloud_security_group" "default" {
  name        = "web accessibility"
  description = "make it accessible for both production and stage ports"
}

# Create security group rule allow web request
resource "tencentcloud_security_group_rule" "web" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}

# Create security group rule allow ssh request
resource "tencentcloud_security_group_rule" "ssh" {
  security_group_id = tencentcloud_security_group.default.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
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
- Assume role

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `secret_id` `secret_key` and `region` in-line in the tencentcloud provider block:

Usage:

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"
}
```

### Environment variables

You can provide your credentials via `TENCENTCLOUD_SECRET_ID` and `TENCENTCLOUD_SECRET_KEY` environment variables,
representing your TencentCloud Secret Id and Secret Key respectively. `TENCENTCLOUD_REGION` is also used, if applicable:

```hcl
provider "tencentcloud" {
}
```

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="my-secret-id"
$ export TENCENTCLOUD_SECRET_KEY="my-secret-key"
$ export TENCENTCLOUD_REGION="ap-guangzhou"
$ terraform plan
```

### Assume role

If provided with an assume role, Terraform will attempt to assume this role using the supplied credentials. Assume role can be provided by adding an `assume_role_arn`, `assume_role_session_name`, `assume_role_session_duration` and `assume_role_policy`(optional) in-line in the tencentcloud provider block:

Usage:

```hcl
provider "tencentcloud" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
  region     = "ap-guangzhou"

  assume_role {
    role_arn         = "my-role-arn"
    session_name     = "my-session-name"
    policy           = "my-role-policy"
    session_duration = 3600
  }
}
```

The `assume_role_arn`, `assume_role_session_name`, `assume_role_session_duration` can also provided via `TENCENTCLOUD_ASSUME_ROLE_ARN`, `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME` and `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION` environment variables.

Usage:

```shell
$ export TENCENTCLOUD_SECRET_ID="my-secret-id"
$ export TENCENTCLOUD_SECRET_KEY="my-secret-key"
$ export TENCENTCLOUD_REGION="ap-guangzhou"
$ export TENCENTCLOUD_ASSUME_ROLE_ARN="my-role-arn"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME="my-session-name"
$ export TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION=3600
$ terraform plan
```

## Argument Reference

In addition to generic provider arguments (e.g. alias and version), the following arguments are supported in the TencentCloud provider block:

* `secret_id` - (Required) This is the TencentCloud secret id. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.
* `secret_key` - (Required) This is the TencentCloud secret key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.
* `security_token` - (Optional) TencentCloud security token of temporary access credentials. It can also be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).
* `region` - (Required) This is the TencentCloud region. It must be provided, but it can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is `ap-guangzhou`.
* `assume_role` - (Optional, Available in 1.33.1+) An `assume_role` block (documented below). If provided, terraform will attempt to assume this role using the supplied credentials. Only one `assume_role` block may be in the configuration.
* `protocol` - (Optional, Available in 1.37.0+) The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.
* `domain` - (Optional, Available in 1.37.0+) The root domain of the API request, Default is `tencentcloudapi.com`.
The nested `assume_role` block supports the following:
* `role_arn` - (Required) The ARN of the role to assume. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN` environment variable.
* `session_name` - (Required) The session name to use when making the AssumeRole call. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME` environment variable.
* `session_duration` - (Required) The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can also be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION` environment variable.
* `policy` - (Optional) A more restrictive policy to apply to the temporary credentials. This gives you a way to further restrict the permissions for the resulting temporary security credentials. You cannot use the passed policy to grant permissions that are in excess of those allowed by the access policy of the role that is being assumed.
