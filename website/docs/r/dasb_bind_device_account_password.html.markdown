---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_bind_device_account_password"
sidebar_current: "docs-tencentcloud-resource-dasb_bind_device_account_password"
description: |-
  Provides a resource to create a dasb bind_device_account_password
---

# tencentcloud_dasb_bind_device_account_password

Provides a resource to create a dasb bind_device_account_password

## Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_account_password" "example" {
  device_account_id = 16
  password          = "TerraformPassword"
}
```

## Argument Reference

The following arguments are supported:

* `device_account_id` - (Required, Int, ForceNew) Host account ID.
* `password` - (Required, String, ForceNew) Host account password.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



