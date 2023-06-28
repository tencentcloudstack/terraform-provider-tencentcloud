---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_security_group_attachment"
sidebar_current: "docs-tencentcloud-resource-redis_security_group_attachment"
description: |-
  Provides a resource to create a redis security_group_attachment
---

# tencentcloud_redis_security_group_attachment

Provides a resource to create a redis security_group_attachment

## Example Usage

```hcl
resource "tencentcloud_redis_security_group_attachment" "security_group_attachment" {
  instance_id       = "crs-jf4ico4v"
  security_group_id = "sg-edmur627"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `security_group_id` - (Required, String, ForceNew) Security group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis security_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_redis_security_group_attachment.security_group_attachment instance_id#security_group_id
```

