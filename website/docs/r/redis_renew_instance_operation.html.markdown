---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_renew_instance_operation"
sidebar_current: "docs-tencentcloud-resource-redis_renew_instance_operation"
description: |-
  Provides a resource to create a redis renew_instance_operation
---

# tencentcloud_redis_renew_instance_operation

Provides a resource to create a redis renew_instance_operation

## Example Usage

```hcl
resource "tencentcloud_redis_renew_instance_operation" "renew_instance_operation" {
  instance_id     = "crs-c1nl9rpv"
  period          = 1
  modify_pay_mode = "prepaid"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `period` - (Required, Int, ForceNew) Purchase duration, in months.
* `modify_pay_mode` - (Optional, String, ForceNew) Identifies whether the billing model is modified:The current instance billing mode is pay-as-you-go, which is prepaid and renewed.The billing mode of the current instance is subscription and you can not set this parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



