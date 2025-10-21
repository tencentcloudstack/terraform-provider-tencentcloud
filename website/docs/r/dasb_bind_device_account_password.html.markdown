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
resource "tencentcloud_dasb_device" "example" {
  os_name = "Linux"
  ip      = "192.168.0.1"
  port    = 80
  name    = "tf_example"
}

resource "tencentcloud_dasb_device_account" "example" {
  device_id = tencentcloud_dasb_device.example.id
  account   = "root"
}

resource "tencentcloud_dasb_bind_device_account_password" "example" {
  device_account_id = tencentcloud_dasb_device_account.example.id
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



