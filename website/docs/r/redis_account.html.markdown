---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_account"
sidebar_current: "docs-tencentcloud-resource-redis_account"
description: |-
  Provides a resource to create a redis account
---

# tencentcloud_redis_account

Provides a resource to create a redis account

## Example Usage

```hcl
resource "tencentcloud_redis_account" "account" {
  instance_id      = "crs-xxxxxx"
  account_name     = "account_test"
  account_password = "test1234"
  remark           = "master"
  readonly_policy  = ["master"]
  privilege        = "rw"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String) The account name.
* `account_password` - (Required, String) 1: Length 8-30 digits, it is recommended to use a password of more than 12 digits; 2: Cannot start with `/`; 3: Include at least two items: a.Lowercase letters `a-z`; b.Uppercase letters `A-Z` c.Numbers `0-9`;  d.`()`~!@#$%^&*-+=_|{}[]:;<>,.?/`.
* `instance_id` - (Required, String) The ID of instance.
* `privilege` - (Required, String) Read and write policy: Enter R and RW to indicate read-only, read-write, cannot be empty when modifying operations.
* `readonly_policy` - (Required, Set: [`String`]) Routing policy: Enter master or replication, which indicates the master node or slave node, cannot be empty when modifying operations.
* `remark` - (Optional, String) Remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis account can be imported using the id, e.g.

```
terraform import tencentcloud_redis_account.account crs-xxxxxx#account_test
```

