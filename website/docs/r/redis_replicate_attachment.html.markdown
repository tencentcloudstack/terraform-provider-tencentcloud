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
  group_id           = "crs-rpl-c1nl9rpv"
  master_instance_id = "crs-c1nl9rpv"
  instance_ids       = []
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) The ID of group.
* `instance_ids` - (Required, Set: [`String`]) All instance ids of the replication group.
* `master_instance_id` - (Required, String) The ID of master instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis replicate_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_redis_replicate_attachment.replicate_attachment replicate_attachment_id
```

