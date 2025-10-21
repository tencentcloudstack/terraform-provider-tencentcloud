---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_restart_instance"
sidebar_current: "docs-tencentcloud-resource-cdwpg_restart_instance"
description: |-
  Provides a resource to create a cdwpg cdwpg_restart_instance
---

# tencentcloud_cdwpg_restart_instance

Provides a resource to create a cdwpg cdwpg_restart_instance

## Example Usage

```hcl
resource "tencentcloud_cdwpg_restart_instance" "cdwpg_restart_instance" {
  instance_id = "cdwpg-zpiemnyd"
  node_types  = ["gtm"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id (e.g., "cdwpg-xxxx").
* `node_ids` - (Optional, Set: [`String`], ForceNew) Node ids to restart (specify nodes to reboot).
* `node_types` - (Optional, Set: [`String`], ForceNew) Node types to restart (gtm/cn/dn).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



