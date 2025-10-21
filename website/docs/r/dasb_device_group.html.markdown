---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_device_group"
sidebar_current: "docs-tencentcloud-resource-dasb_device_group"
description: |-
  Provides a resource to create a dasb device_group
---

# tencentcloud_dasb_device_group

Provides a resource to create a dasb device_group

## Example Usage

```hcl
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example"
  department_id = "1.2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Device group name, the maximum length is 32 characters.
* `department_id` - (Optional, String) The ID of the department to which the asset group belongs, such as: 1.2.3 name, with a maximum length of 32 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb device_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_group.example 36
```

