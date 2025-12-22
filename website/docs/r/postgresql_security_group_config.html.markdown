---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_security_group_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_security_group_config"
description: |-
  Provides a resource to create a postgresql security_group_config
---

# tencentcloud_postgresql_security_group_config

Provides a resource to create a postgresql security_group_config

## Example Usage

Set security group for the sepcified postgres instance

```hcl
resource "tencentcloud_postgresql_security_group_config" "security_group_config" {
  security_group_id_set = [local.sg_id, local.sg_id2]
  db_instance_id        = local.pgsql_id
}
```

Set security group for the specified readonly group

```hcl
resource "tencentcloud_postgresql_readonly_group" "group" {
  master_db_instance_id       = local.pgsql_id
  name                        = "tf_test_ro_sg"
  project_id                  = 0
  subnet_id                   = local.subnet_id
  vpc_id                      = local.vpc_id
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

resource "tencentcloud_postgresql_security_group_config" "security_group_config" {
  security_group_id_set = [local.sg_id, local.sg_id2]
  read_only_group_id    = tencentcloud_postgresql_readonly_group.group.id
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id_set` - (Required, Set: [`String`]) Information of security groups in array.
* `db_instance_id` - (Optional, String) Instance ID. Either this parameter or ReadOnlyGroupId must be passed in. If both parameters are passed in, ReadOnlyGroupId will be ignored.
* `read_only_group_id` - (Optional, String) RO group ID. Either this parameter or DBInstanceId must be passed in. To query the security groups associated with the RO groups, only pass in ReadOnlyGroupId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



