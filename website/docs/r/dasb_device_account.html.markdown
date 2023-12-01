---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_device_account"
sidebar_current: "docs-tencentcloud-resource-dasb_device_account"
description: |-
  Provides a resource to create a dasb device_account
---

# tencentcloud_dasb_device_account

Provides a resource to create a dasb device_account

## Example Usage

```hcl
resource "tencentcloud_dasb_device_account" "example" {
  device_id = 100
  account   = "root"
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required, String, ForceNew) Device account.
* `device_id` - (Required, Int, ForceNew) Device ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb device_account can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_account.example 11
```

