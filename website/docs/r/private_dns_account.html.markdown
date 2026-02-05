---
subcategory: "Private DNS"
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

### Use with Private DNS Zone

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}

resource "tencentcloud_private_dns_zone" "example" {
  domain = "example.com"

  account_vpc_set {
    uniq_vpc_id = "vpc-xxxxx"
    region      = "ap-guangzhou"
    uin         = tencentcloud_private_dns_account.example.account_uin
  }
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

* `account_uin` - (Required, String, ForceNew) Uin of the associated account. Changing this will recreate the account association.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, which is the Uin of the associated account.
* `account` - Email of the associated account.
* `nickname` - Nickname of the associated account.

## Import

Private DNS account association can be imported using the account Uin, e.g.

```
$ terraform import tencentcloud_private_dns_account.example 100123456789
```
