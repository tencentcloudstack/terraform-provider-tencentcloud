---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table_entry_config"
sidebar_current: "docs-tencentcloud-resource-route_table_entry_config"
description: |-
  Provides a resource to create a vpc route table entry config
---

# tencentcloud_route_table_entry_config

Provides a resource to create a vpc route table entry config

~> **NOTE:** When setting the route item switch, do not use it together with resource `tencentcloud_route_table_entry`.

## Example Usage

### Enable route item

```hcl
resource "tencentcloud_route_table_entry_config" "example" {
  route_table_id = "rtb-8425lgjy"
  route_item_id  = "rti-4f6efqwn"
  disabled       = false
}
```

### Disable route item

```hcl
resource "tencentcloud_route_table_entry_config" "example" {
  route_table_id = "rtb-8425lgjy"
  route_item_id  = "rti-4f6efqwn"
  disabled       = true
}
```

## Argument Reference

The following arguments are supported:

* `disabled` - (Required, Bool) Whether the entry is disabled.
* `route_item_id` - (Required, String, ForceNew) ID of route table entry.
* `route_table_id` - (Required, String, ForceNew) Route table ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc route table entry config can be imported using the id, e.g.

```
terraform import tencentcloud_route_table_entry_config.example rtb-8425lgjy#rti-4f6efqwn
```

