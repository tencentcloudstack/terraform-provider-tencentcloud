---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_account"
sidebar_current: "docs-tencentcloud-resource-private_dns_account"
description: |-
  Provides a resource to create a privatedns account
---

# tencentcloud_private_dns_account

Provides a resource to create a privatedns account

## Example Usage

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account {
    uin = "100022770160"
  }
}
```

### Or

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account {
    uin      = "100022770160"
    account  = "example@tencent.com"
    nickname = "tf-example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required, List, ForceNew) Private DNS account.

The `account` object supports the following:

* `uin` - (Required, String, ForceNew) Root account UIN.
* `account` - (Optional, String, ForceNew) Root account name.
* `nickname` - (Optional, String, ForceNew) Account name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

privatedns account can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_account.example 100022770160
```

