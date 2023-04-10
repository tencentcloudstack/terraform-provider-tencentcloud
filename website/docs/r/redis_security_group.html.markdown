---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_security_group"
sidebar_current: "docs-tencentcloud-resource-redis_security_group"
description: |-
  Provides a resource to create a redis security_group
---

# tencentcloud_redis_security_group

Provides a resource to create a redis security_group

## Example Usage

```hcl
resource "tencentcloud_redis_security_group" "security_group" {
  instance_id       = "crs-c1nl9rpv"
  security_group_id = "sg-cyules4s5"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `security_group_id` - (Required, String, ForceNew) Security group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis security_group can be imported using the instance_id#security_group_id, e.g.

```
terraform import tencentcloud_redis_security_group.security_group crs-c1nl9rpv#sg-cyules4s5
```

