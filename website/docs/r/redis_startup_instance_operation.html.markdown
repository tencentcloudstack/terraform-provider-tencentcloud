---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_startup_instance_operation"
sidebar_current: "docs-tencentcloud-resource-redis_startup_instance_operation"
description: |-
  Provides a resource to create a redis startup instance operation
---

# tencentcloud_redis_startup_instance_operation

Provides a resource to create a redis startup instance operation

## Example Usage

### Recover the redis instance that has been isolated

```hcl
resource "tencentcloud_redis_startup_instance_operation" "example" {
  instance_id = "crs-c1nl9rpv"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



