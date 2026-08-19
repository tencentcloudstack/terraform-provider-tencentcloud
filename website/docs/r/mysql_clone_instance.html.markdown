---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_clone_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_clone_instance"
description: |-
  Provides a resource to create a TencentDB for MySQL (CDB) clone instance from an existing source instance, optionally rolling back to a specified time point or backup set.
---

# tencentcloud_mysql_clone_instance

Provides a resource to create a TencentDB for MySQL (CDB) clone instance from an existing source instance, optionally rolling back to a specified time point or backup set.

## Example Usage

### Clone instance with rollback time

```hcl
resource "tencentcloud_mysql_clone_instance" "example" {
  instance_id             = "cdb-c1nl9rpv"
  specified_rollback_time = "2024-01-01 12:00:00"
  memory                  = 4000
  volume                  = 200
  cpu                     = 2
  instance_name           = "tf-clone-example"
  uniq_vpc_id             = "vpc-i5yyodl9"
  uniq_subnet_id          = "subnet-hhi88a58"
  protect_mode            = 0
  deploy_mode             = 0
  slave_zone              = "ap-guangzhou-3"
  device_type             = "UNIVERSAL"
  project_id              = 0
  pay_type                = "USED_PAID"
  zone                    = "ap-guangzhou-3"
}
```

### Clone instance with backup id

```hcl
resource "tencentcloud_mysql_clone_instance" "example" {
  instance_id         = "cdb-c1nl9rpv"
  specified_backup_id = 1000001
  memory              = 4000
  volume              = 200
  cpu                 = 2
  instance_name       = "tf-clone-example"
  uniq_vpc_id         = "vpc-i5yyodl9"
  uniq_subnet_id      = "subnet-hhi88a58"
  protect_mode        = 0
  deploy_mode         = 0
  slave_zone          = "ap-guangzhou-3"
  device_type         = "UNIVERSAL"
  project_id          = 0
  pay_type            = "USED_PAID"
  zone                = "ap-guangzhou-3"
}
```

### Clone instance with resource tags

```hcl
resource "tencentcloud_mysql_clone_instance" "example" {
  instance_id             = "cdb-c1nl9rpv"
  specified_rollback_time = "2024-01-01 12:00:00"
  memory                  = 4000
  volume                  = 200
  cpu                     = 2
  instance_name           = "tf-clone-example"
  uniq_vpc_id             = "vpc-i5yyodl9"
  uniq_subnet_id          = "subnet-hhi88a58"
  protect_mode            = 0
  deploy_mode             = 0
  slave_zone              = "ap-guangzhou-3"
  device_type             = "UNIVERSAL"
  project_id              = 0
  pay_type                = "USED_PAID"
  zone                    = "ap-guangzhou-3"

  resource_tags {
    key   = "createBy"
    value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Source MySQL (CDB) instance ID to clone from.
* `backup_zone` - (Optional, String) Slave 2 zone. Updatable via `UpgradeDBInstance`.
* `cage_id` - (Optional, String, ForceNew) Financial cage ID.
* `cluster_topology` - (Optional, List) Cloud disk node topology. Updatable via `UpgradeDBInstance`.
* `cpu` - (Optional, Int) CPU cores. Updatable via `UpgradeDBInstance`.
* `deploy_group_id` - (Optional, String, ForceNew) Placement group ID.
* `deploy_mode` - (Optional, Int) Deploy mode. 0 - single zone, 1 - multi zone. Updatable via `UpgradeDBInstance`.
* `device_type` - (Optional, String) Instance type. Updatable via `UpgradeDBInstance`.
* `dry_run` - (Optional, Bool, ForceNew) Dry-run flag. true: only pre-check, false: send normal request.
* `fourth_zone` - (Optional, String) Slave 3 zone. Updatable via `UpgradeDBInstance`.
* `instance_name` - (Optional, String, ForceNew) Cloned instance name.
* `instance_nodes` - (Optional, Int, ForceNew) Instance node count.
* `master_zone` - (Optional, String, ForceNew) Master zone.
* `memory` - (Optional, Int) Instance memory in MB. Updatable via `UpgradeDBInstance`.
* `pay_type` - (Optional, String, ForceNew) Payment type. PRE_PAID - prepaid, USED_PAID - postpaid. Default is postpaid.
* `period` - (Optional, Int, ForceNew) Instance duration in months. Required when `pay_type` is PRE_PAID.
* `project_id` - (Optional, Int, ForceNew) Project ID. Default is 0.
* `protect_mode` - (Optional, Int) Data replication mode. 0 - async, 1 - semisync, 2 - strongsync. Updatable via `UpgradeDBInstance`.
* `resource_tags` - (Optional, List, ForceNew) Instance tags.
* `security_group` - (Optional, List: [`String`], ForceNew) Security group IDs.
* `slave_zone` - (Optional, String) Slave 1 zone. Updatable via `UpgradeDBInstance`.
* `specified_backup_id` - (Optional, Int, ForceNew) Backup file ID to clone from. Mutually exclusive with `specified_rollback_time`.
* `specified_rollback_time` - (Optional, String, ForceNew) Rollback time (yyyy-mm-dd hh:mm:ss). Mutually exclusive with `specified_backup_id`.
* `specified_sub_backup_id` - (Optional, Int, ForceNew) Cross-region backup ID.
* `src_region` - (Optional, String, ForceNew) Source instance region for cross-region clone.
* `uniq_subnet_id` - (Optional, String, ForceNew) Subnet ID. If `uniq_vpc_id` is set, this value is required.
* `uniq_vpc_id` - (Optional, String, ForceNew) VPC ID.
* `volume` - (Optional, Int) Instance disk size in GB. Updatable via `UpgradeDBInstance`.
* `zone` - (Optional, String, ForceNew) Instance zone.

The `cluster_topology` object supports the following:

* `read_only_nodes` - (Optional, List) RO node topology list.
* `read_write_node` - (Optional, List) RW node topology.

The `read_only_nodes` object of `cluster_topology` supports the following:

* `is_random_zone` - (Optional, String) Whether distributed in a random zone. YES - random zone.
* `node_id` - (Optional, String) Node ID.
* `zone` - (Optional, String) RO node zone.

The `read_write_node` object of `cluster_topology` supports the following:

* `node_id` - (Optional, String) Node ID.
* `zone` - (Optional, String) RW node zone.

The `resource_tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `async_request_id` - Async request ID returned by Create/Update APIs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `1h0m`) Used when creating the resource.
* `update` - (Defaults to `1h0m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

TencentDB for MySQL (CDB) clone instance can be imported using the cloned instance id, e.g.

```
terraform import tencentcloud_mysql_clone_instance.example cdb-bcet7sdb
```

