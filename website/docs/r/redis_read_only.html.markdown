---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_read_only"
sidebar_current: "docs-tencentcloud-resource-redis_read_only"
description: |-
  Provides a resource to create a redis read_only
---

# tencentcloud_redis_read_only

Provides a resource to create a redis read_only

## Example Usage

```hcl
resource "tencentcloud_redis_read_only" "read_only" {
  instance_id = "crs-c1nl9rpv"
  input_mode  = "0"
}
```

## Argument Reference

The following arguments are supported:

* `input_mode` - (Required, String) Instance input mode: `0`: read-write; `1`: read-only.
* `instance_id` - (Required, String) The ID of instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis read_only can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_read_only.read_only crs-c1nl9rpv
```

