---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_readonly_instance"
description: |-
  Provides a SQL Server instance resource to create read-only database instances.
---

# tencentcloud_sqlserver_readonly_instance

Provides a SQL Server instance resource to create read-only database instances.

## Example Usage

```hcl
resource "tencentcloud_mysql_readonly_instance" "foo" {
  name                = "tf_sqlserver_instance_ro"
  availability_zone   = "ap-guangzhou-4"
  charge_type         = "POSTPAID_BY_HOUR"
  vpc_id              = "` + defaultVpcId + `"
  subnet_id           = "` + defaultSubnetId + `"
  memory              = 2
  storage             = 10
  master_instance_id  = tencentcloud_sqlserver_instance.test.id
  readonly_group_type = 1
  force_upgrade       = true
}
```

## Argument Reference

The following arguments are supported:

* `master_instance_id` - (Required, ForceNew) Indicates the master instance ID of recovery instances.
* `memory` - (Required) Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
* `name` - (Required) Name of the SQL Server instance.
* `readonly_group_type` - (Required, ForceNew) Type of readonly group. Valid values: 1, 3. 1 for one auto-assigned readonly instance per one readonly group, 2 for creating new readonly group, 3 for all exist readonly instances stay in the exist readonly group. For now, only 1 and 3 are supported.
* `storage` - (Required) Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
* `availability_zone` - (Optional, ForceNew) Availability zone.
* `charge_type` - (Optional, ForceNew) Pay type of the SQL Server instance. For now, only `POSTPAID_BY_HOUR` is valid.
* `force_upgrade` - (Optional, ForceNew) Indicate that the master instance upgrade or not. `true` for upgrading the master SQL Server instance to cluster type by force. Default is false. Note: this is not supported with `DUAL`(ha_type), `2017`(engine_version) master SQL Server instance, for it will cause ha_type of the master SQL Server instance change.
* `readonly_group_id` - (Optional) ID of the readonly group that this instance belongs to. When `readonly_group_type` set value `3`, it must be set with valid value.
* `security_groups` - (Optional) Security group bound to the instance.
* `subnet_id` - (Optional, ForceNew) ID of subnet.
* `tags` - (Optional) The tags of the SQL Server.
* `vpc_id` - (Optional, ForceNew) ID of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server instance.
* `ro_flag` - Readonly flag. `RO` for readonly instance, `MASTER` for master instance,  `` for not readonly instance.
* `status` - Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
* `vip` - IP for private access.
* `vport` - Port for private access.


## Import

SQL Server readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_readonly_instance.foo mssqlro-3cdq7kx5
```

