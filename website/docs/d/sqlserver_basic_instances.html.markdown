---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_basic_instances"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_basic_instances"
description: |-
  Use this data source to query SQL Server basic instances
---

# tencentcloud_sqlserver_basic_instances

Use this data source to query SQL Server basic instances

## Example Usage

```hcl
resource "tencentcloud_sqlserver_basic_instance" "test" {
  name              = "tf_sqlserver_basic_instance"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-26w7r56z"
  subnet_id         = "subnet-lvlr6eeu"
  machine_type      = "CLOUD_PREMIUM"
  project_id        = 0
  memory            = 2
  storage           = 10
  cpu               = 1
  security_groups   = ["sg-nltpbqg1"]

  tags = {
    "test" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the SQL Server basic instance to be query.
* `project_id` - (Optional) Project ID of the SQL Server basic instance to be query.
* `result_output_file` - (Optional) Used to save results.
* `subnet_id` - (Optional) Subnet ID of the SQL Server basic instance to be query.
* `vpc_id` - (Optional) Vpc ID of the SQL Server basic instance to be query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of SQL Server basic instances. Each element contains the following attributes.
  * `availability_zone` - Availability zone.
  * `charge_type` - Pay type of the SQL Server basic instance. For now, only `POSTPAID_BY_HOUR` is valid.
  * `cpu` - The CPU number of the SQL Server basic instance.
  * `create_time` - Create time of the SQL Server basic instance.
  * `engine_version` - Version of the SQL Server basic database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
  * `id` - ID of the SQL Server basic instance.
  * `memory` - Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
  * `name` - Name of the SQL Server basic instance.
  * `project_id` - Project ID, default value is 0.
  * `status` - Status of the SQL Server basic instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
  * `storage` - Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
  * `subnet_id` - ID of subnet.
  * `tags` - Tags of the SQL Server basic instance.
  * `used_storage` - Used storage.
  * `vip` - IP for private access.
  * `vpc_id` - ID of VPC.
  * `vport` - Port for private access.


