---
subcategory: "tdcpg"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdcpg_instance"
sidebar_current: "docs-tencentcloud-resource-tdcpg_instance"
description: |-
  Provides a resource to create a tdcpg instance
---

# tencentcloud_tdcpg_instance

Provides a resource to create a tdcpg instance

## Example Usage

```hcl
resource "tencentcloud_tdcpg_instance" "instance1" {
  cluster_id    = "cluster_id"
  cpu           = 1
  memory        = 1
  instance_name = "instance_name"
}

resource "tencentcloud_tdcpg_instance" "instance2" {
  cluster_id       = "cluster_id"
  cpu              = 1
  memory           = 2
  instance_name    = "instance_name"
  operation_timing = "IMMEDIATE"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) cluster id.
* `cpu` - (Required, Int) cpu cores.
* `memory` - (Required, Int) memory size.
* `instance_name` - (Optional, String) instance name.
* `operation_timing` - (Optional, String) operation timing, optional value is IMMEDIATE or MAINTAIN_PERIOD.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdcpg instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_instance.instance cluster_id#instance_id
```

