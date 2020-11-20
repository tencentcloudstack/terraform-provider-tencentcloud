---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_basic_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_basic_instance"
description: |-
  Provides a SQL Server instance resource to create basic database instances.
---

# tencentcloud_sqlserver_basic_instance

Provides a SQL Server instance resource to create basic database instances.

## Example Usage

```hcl
resource "tencentcloud_sqlserver_basic_instance" "foo" {
  name                   = "example"
  availability_zone      = var.availability_zone
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = "vpc-26w7r56z"
  subnet_id              = "subnet-lvlr6eeu"
  project_id             = 0
  memory                 = 2
  storage                = 20
  cpu                    = 1
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = ["sg-nltpbqg1"]

  tags = {
    "test" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required) The CPU number of the SQL Server basic instance.
* `machine_type` - (Required) The host type of the purchased instance, `CLOUD_PREMIUM` for virtual machine high-performance cloud disk, `CLOUD_SSD` for virtual machine SSD cloud disk.
* `memory` - (Required) Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
* `name` - (Required) Name of the SQL Server basic instance.
* `storage` - (Required) Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
* `auto_renew` - (Optional) Automatic renewal sign. 0 for normal renewal, 1 for automatic renewal, the default is 1 automatic renewal. Only valid when purchasing a prepaid instance.
* `auto_voucher` - (Optional) Whether to use the voucher automatically; 1 for yes, 0 for no, the default is 0.
* `availability_zone` - (Optional, ForceNew) Availability zone.
* `charge_type` - (Optional, ForceNew) Pay type of the SQL Server basic instance. For now, only `POSTPAID_BY_HOUR` is valid.
* `engine_version` - (Optional, ForceNew) Version of the SQL Server basic database engine. Allowed values are `2008R2`(SQL Server 2008 Enerprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
* `maintenance_start_time` - (Optional) Start time of the maintenance in one day, format like `HH:mm`.
* `maintenance_time_span` - (Optional) The timespan of maintenance in one day, unit is hour.
* `maintenance_week_set` - (Optional) A list of integer indicates weekly maintenance. For example, [1,7] presents do weekly maintenance on every Monday and Sunday.
* `period` - (Optional) Purchase instance period, the default value is 1, which means one month. The value does not exceed 48.
* `project_id` - (Optional) Project ID, default value is 0.
* `security_groups` - (Optional) Security group bound to the instance.
* `subnet_id` - (Optional, ForceNew) ID of subnet.
* `tags` - (Optional) The tags of the SQL Server basic instance.
* `voucher_ids` - (Optional) An array of voucher IDs, currently only one can be used for a single order.
* `vpc_id` - (Optional, ForceNew) ID of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server basic instance.
* `status` - Status of the SQL Server basic instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
* `vip` - IP for private access.
* `vport` - Port for private access.


## Import

SQL Server basic instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_basic_instance.foo mssql-3cdq7kx5
```

