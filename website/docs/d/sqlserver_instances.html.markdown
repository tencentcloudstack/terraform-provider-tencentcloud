---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instances"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_instances"
description: |-
  Use this data source to query SQL Server instances
---

# tencentcloud_sqlserver_instances

Use this data source to query SQL Server instances

## Example Usage

### Filter instance by Id

```hcl
data "tencentcloud_sqlserver_instances" "example_id" {
  id = "mssql-3l3fgqn7"
}
```

### Filter instance by project Id

```hcl
data "tencentcloud_sqlserver_instances" "example_project" {
  project_id = 0
}
```

### Filter instance by VPC/Subnet

```hcl
data "tencentcloud_sqlserver_instances" "example_vpc" {
  vpc_id    = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the SQL Server instance to be query.
* `name` - (Optional, String) Name of the SQL Server instance to be query.
* `project_id` - (Optional, Int) Project ID of the SQL Server instance to be query.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) Subnet ID of the SQL Server instance to be query.
* `vpc_id` - (Optional, String) Vpc ID of the SQL Server instance to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of SQL Server instances. Each element contains the following attributes.
  * `availability_zone` - Availability zone.
  * `charge_type` - Pay type of the SQL Server instance. For now, only `POSTPAID_BY_HOUR` is valid.
  * `create_time` - Create time of the SQL Server instance.
  * `engine_version` - Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
  * `ha_type` - Instance type. `DUAL` (dual-server high availability), `CLUSTER` (cluster).
  * `id` - ID of the SQL Server instance.
  * `memory` - Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
  * `name` - Name of the SQL Server instance.
  * `project_id` - Project ID, default value is 0.
  * `ro_flag` - Readonly flag. `RO` (read-only instance), `MASTER` (primary instance with read-only instances). If it is left empty, it refers to an instance which is not read-only and has no RO group.
  * `status` - Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
  * `storage` - Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
  * `subnet_id` - ID of subnet.
  * `tags` - Tags of the SQL Server instance.
  * `used_storage` - Used storage.
  * `vip` - IP for private access.
  * `vpc_id` - ID of VPC.
  * `vport` - Port for private access.


