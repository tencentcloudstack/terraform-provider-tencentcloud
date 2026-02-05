---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_account"
sidebar_current: "docs-tencentcloud-resource-private_dns_account"
description: |-
  Provides a resource to create a Private DNS account association.
---

# tencentcloud_private_dns_account

Provides a resource to create a Private DNS account association.

This resource is used to associate an account with Private DNS, enabling cross-account VPC binding for private zones.

~> **NOTE:** Once an account is associated, it can be used to bind VPCs from that account to private DNS zones.

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}
```

### Output Account Information

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}

output "account_email" {
  value = tencentcloud_private_dns_account.example.account
}

output "account_nickname" {
  value = tencentcloud_private_dns_account.example.nickname
}
```

## Argument Reference

The following arguments are supported:

* `account_uin` - (Required, String, ForceNew) Uin of the associated account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account` - Email of the associated account.
* `nickname` - Nickname of the associated account.


