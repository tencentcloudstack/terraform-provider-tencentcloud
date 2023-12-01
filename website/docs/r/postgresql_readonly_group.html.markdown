---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_group"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_group"
description: |-
  Use this resource to create postgresql readonly group.
---

# tencentcloud_postgresql_readonly_group

Use this resource to create postgresql readonly group.

## Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group" "group" {
  master_db_instance_id       = "postgres-gzg9jb2n"
  name                        = "world"
  project_id                  = 0
  vpc_id                      = "vpc-86v957zb"
  subnet_id                   = "subnet-enm92y0m"
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
  #  security_groups_ids = []
}
```

## Argument Reference

The following arguments are supported:

* `master_db_instance_id` - (Required, String, ForceNew) Primary instance ID.
* `max_replay_lag` - (Required, Int) Delay threshold in ms.
* `max_replay_latency` - (Required, Int) Delayed log size threshold in MB.
* `min_delay_eliminate_reserve` - (Required, Int) The minimum number of read-only replicas that must be retained in an RO group.
* `name` - (Required, String) RO group name.
* `project_id` - (Required, Int) Project ID.
* `replay_lag_eliminate` - (Required, Int) Whether to remove a read-only replica from an RO group if the delay between the read-only replica and the primary instance exceeds the threshold. Valid values: 0 (no), 1 (yes).
* `replay_latency_eliminate` - (Required, Int) Whether to remove a read-only replica from an RO group if the sync log size difference between the read-only replica and the primary instance exceeds the threshold. Valid values: 0 (no), 1 (yes).
* `subnet_id` - (Required, String) VPC subnet ID.
* `vpc_id` - (Required, String) VPC ID.
* `security_groups_ids` - (Optional, Set: [`String`]) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the postgresql instance.


