---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_cpu_scale"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_cpu_scale"
description: |-
  Provides a resource to manage MongoDB instance CPU elastic scaling configuration.
---

# tencentcloud_mongodb_instance_cpu_scale

Provides a resource to manage MongoDB instance CPU elastic scaling configuration.

~> **NOTE:** This resource is used to enable or disable CPU elastic scaling for MongoDB instances. When the resource is destroyed, CPU elastic scaling will be disabled automatically.

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_cpu_scale" "example" {
  instance_id = "cmgo-xxxxxxxx"
  extra_cpu   = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) MongoDB instance ID, for example: cmgo-xxxxxxxx.
* `extra_cpu` - (Optional, Int) Extra CPU cores for elastic scaling. Each node will be scaled up by this number of cores. Setting this value enables CPU elastic scaling. Removing this field or setting it to 0 disables CPU elastic scaling.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `flow_id` - Task flow ID returned by the last operation.


## Import

MongoDB instance CPU elastic scaling configuration can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_instance_cpu_scale.example cmgo-xxxxxxxx
```

