---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_replica_readonly"
sidebar_current: "docs-tencentcloud-resource-redis_replica_readonly"
description: |-
  Provides a resource to create a redis replica_readonly
---

# tencentcloud_redis_replica_readonly

Provides a resource to create a redis replica_readonly

## Example Usage

```hcl
resource "tencentcloud_redis_replica_readonly" "replica_readonly" {
  instance_id     = "crs-c1nl9rpv"
  readonly_policy = ["master"]
  operate         = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `operate` - (Required, String) The replica is read-only, `enable` - enable read-write splitting, `disable`- disable read-write splitting.
* `readonly_policy` - (Optional, Set: [`String`]) Routing policy: Enter `master` or `replication`, which indicates the master node or slave node.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



