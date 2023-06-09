---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_replicate_attachment"
sidebar_current: "docs-tencentcloud-resource-redis_replicate_attachment"
description: |-
  Provides a resource to create a redis replicate_attachment
---

# tencentcloud_redis_replicate_attachment

Provides a resource to create a redis replicate_attachment

## Example Usage

```hcl
resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
  instance_id   = "crs-c1nl9rpv"
  group_id      = "crs-rpl-c1nl9rpv"
  instance_role = "rw"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) The ID of group.
* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `instance_role` - (Required, String, ForceNew) Assign roles to instances added to the replication group.:rw: read-write.r: read-only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis replicate_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_redis_replicate_attachment.replicate_attachment replicate_attachment_id
```

