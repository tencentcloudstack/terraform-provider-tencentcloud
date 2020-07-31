---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instances"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_instances"
description: |-
  Use this data source to query SQL Server instances
---

# tencentcloud_sqlserver_instances

Use this data source to query SQL Server instances

## Example Usage

```hcl
data "tencentcloud_sqlserver_instances" "vpc" {
  vpc_id    = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
}

data "tencentcloud_sqlserver_instances" "project" {
  project_id = 0
}

data "tencentcloud_sqlserver_instances" "id" {
  id = "postgres-h9t4fde1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the SQL Server instance to be query.
* `project_id` - (Optional) Project ID of the SQL Server instance to be query.
* `result_output_file` - (Optional) Used to save results.
* `subnet_id` - (Optional) Subnet ID of the SQL Server instance to be query.
* `vpc_id` - (Optional) Vpc ID of the SQL Server instance to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of SQL Server instances. Each element contains the following attributes.
  * `availability_zone` - Availability zone.
  * `charge_type` - Pay type of the SQL Server instance. For now, only `POSTPAID_BY_HOUR` is valid.
  * `create_time` - Create time of the SQL Server instance.
  * `engine_version` - Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enerprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
  * `ha_type` - Instance type.
  * `id` - ID of the SQL Server instance.
  * `memory` - Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
  * `name` - Name of the SQL Server instance.
  * `project_id` - Project ID, default value is 0.
  * `ro_flag` - Readonly flag. `RO` for readonly instance, `MASTER` for master instance,  `` for not readonly instance.
  * `status` - Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
  * `storage` - Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
  * `subnet_id` - ID of subnet.
  * `used_storage` - Used storage.
  * `vip` - IP for private access.
  * `vpc_id` - ID of VPC.
  * `vport` - Port for private access.


