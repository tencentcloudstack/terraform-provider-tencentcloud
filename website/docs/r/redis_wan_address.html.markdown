---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_wan_address"
sidebar_current: "docs-tencentcloud-resource-redis_wan_address"
description: |-
  Provides a resource to create a redis wan_address
---

# tencentcloud_redis_wan_address

Provides a resource to create a redis wan_address

## Example Usage

```hcl
resource "tencentcloud_redis_wan_address" "wan_address" {
  instance_id = "crs-dekqpd8v"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `wan_address` - Allocate Wan Address.


## Import

redis wan_address can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_wan_address.wan_address crs-dekqpd8v
```

