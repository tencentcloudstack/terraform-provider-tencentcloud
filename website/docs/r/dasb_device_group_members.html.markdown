---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_device_group_members"
sidebar_current: "docs-tencentcloud-resource-dasb_device_group_members"
description: |-
  Provides a resource to create a dasb device_group_members
---

# tencentcloud_dasb_device_group_members

Provides a resource to create a dasb device_group_members

## Example Usage

```hcl
resource "tencentcloud_dasb_device" "example" {
  os_name = "Linux"
  ip      = "192.168.0.1"
  port    = 80
  name    = "tf_example"
}

resource "tencentcloud_dasb_device_group" "example" {
  name = "tf_example"
}

resource "tencentcloud_dasb_device_group_members" "example" {
  device_group_id = tencentcloud_dasb_device_group.example.id
  member_id_set   = [tencentcloud_dasb_device.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `device_group_id` - (Required, Int, ForceNew) Device Group ID.
* `member_id_set` - (Required, Set: [`Int`], ForceNew) A collection of device IDs that need to be added to the device group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb device_group_members can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_group_members.example 53#102
```

