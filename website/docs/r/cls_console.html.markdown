---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_console"
sidebar_current: "docs-tencentcloud-resource-cls_console"
description: |-
  Provides a resource to create a CLS DataSight console.
---

# tencentcloud_cls_console

Provides a resource to create a CLS DataSight console.

## Example Usage

### If login_mode is set to 0

```hcl
resource "tencentcloud_cls_console" "example" {
  access_mode     = ["public", "internal"]
  login_mode      = 0
  domain_prefix   = "datasight"
  remarks         = "remarks."
  intranet_type   = 1
  intranet_region = "ap-chongqing"

  accounts {
    user_name  = "your_username1"
    password   = "your_password1"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
    email      = "demo@example.com"
  }

  accounts {
    user_name  = "your_username2"
    password   = "your_password2"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
    email      = "demo@example.com"
  }

  accounts {
    user_name  = "your_username3"
    password   = "your_password3"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
    email      = "demo@example.com"
  }

  menus = [
    "/cls/search",
    "/cls/dashboard",
    "/cls/alarm",
    "/cls/process",
  ]

  access_control_rules {
    access_mode = "public"
    action      = "DENY"
    cidr_blocks = [
      "1.1.1.1",
      "2.2.2.2",
      "3.3.3.3",
    ]
  }

  access_control_rules {
    access_mode = "internal"
    action      = "DENY"
    cidr_blocks = [
      "4.4.4.4",
      "5.5.5.5"
    ]
  }

  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

### If login_mode is set to 1

```hcl
resource "tencentcloud_cls_console" "example" {
  access_mode     = ["internal"]
  login_mode      = 1
  domain_prefix   = "datasight"
  remarks         = "remarks."
  intranet_type   = 1
  intranet_region = "ap-chongqing"

  anonymous_login {
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  menus = [
    "/cls/search",
    "/cls/dashboard",
    "/cls/alarm",
  ]

  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

### If login_mode is set to 2

```hcl
resource "tencentcloud_cls_console" "example2" {
  access_mode     = ["internal"]
  login_mode      = 2
  domain_prefix   = "datasight2"
  remarks         = "remarks."
  intranet_type   = 1
  intranet_region = "ap-chongqing"

  auth_roles {
    role_name  = "role1"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  auth_roles {
    role_name  = "role2"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  auth_roles {
    role_name  = "role3"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  menus = [
    "/cls/search",
    "/cls/dashboard",
    "/cls/alarm",
    "/cls/process",
  ]

  access_control_rules {
    access_mode = "internal"
    action      = "ACCEPT"
    cidr_blocks = [
      "1.1.1.1",
      "2.2.2.2"
    ]
  }

  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_mode` - (Required, List: [`String`]) Access mode list. Valid values: `public` (public network), `internal` (intranet).
* `domain_prefix` - (Required, String) Custom domain prefix.
* `login_mode` - (Required, Int) Login mode. Valid values: `0` (account-password authentication), `1` (anonymous login), `2` (third-party authentication login).
* `access_control_rules` - (Optional, List) Access control rules. Required when `login_mode` is `2` (third-party authentication login) and AccessMode contains `internal` with Action ACCEPT rules.
* `accounts` - (Optional, List) User account information. Required when `login_mode` is `0` (account-password authentication).
* `anonymous_login` - (Optional, List) Anonymous login account information. Required when `login_mode` is `1` (anonymous login).
* `auth_roles` - (Optional, List) Auth role information. Required when `login_mode` is `2` (third-party authentication login).
* `hide_params` - (Optional, List: [`String`]) Custom hidden parameters.
* `intranet_region` - (Optional, String) Intranet region.
* `intranet_type` - (Optional, Int) Intranet type. Default is 0.
* `menus` - (Optional, List: [`String`]) Custom display menus.
* `remarks` - (Optional, String) Remarks.
* `subnet_id` - (Optional, String) Intranet subnet ID.
* `tags` - (Optional, List, ForceNew) Tag bindings. ModifyConsole does not accept tags, so changing this field forces resource replacement.
* `vpc_id` - (Optional, String) Intranet VPC ID.

The `access_control_rules` object supports the following:

* `access_mode` - (Optional, String) Access mode for the rule. Valid values: `public`, `internal`.
* `action` - (Optional, String) Rule action. Valid values: `ACCEPT`, `DROP`.
* `cidr_blocks` - (Optional, List) CIDR blocks or IPs, supporting IPv4 or IPv6.

The `accounts` object supports the following:

* `email` - (Optional, String) Email address used to send verification codes.
* `password` - (Optional, String) User password.
* `secret_id` - (Optional, String) Tencent Cloud account SecretId.
* `secret_key` - (Optional, String) Tencent Cloud account SecretKey.
* `user_name` - (Optional, String) User name.

The `anonymous_login` object supports the following:

* `secret_id` - (Optional, String) Anonymous login account SecretId.
* `secret_key` - (Optional, String) Anonymous login account SecretKey.

The `auth_roles` object supports the following:

* `role_name` - (Optional, String) Auth role name.
* `secret_id` - (Optional, String) SecretId for the auth role permission.
* `secret_key` - (Optional, String) SecretKey for the auth role permission.

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `console_id` - DataSight console ID.
* `domain` - Public network access domain.
* `intranet_domain` - Intranet access domain.


## Import

CLS DataSight console can be imported using the id, e.g.

```
terraform import tencentcloud_cls_console.example clsconsole-0d59e193
```

