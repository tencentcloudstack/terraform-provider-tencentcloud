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
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required, Int) The CPU number of the SQL Server basic instance.
* `machine_type` - (Required, String) The host type of the purchased instance, `CLOUD_PREMIUM` for virtual machine high-performance cloud disk, `CLOUD_SSD` for virtual machine SSD cloud disk, `CLOUD_HSSD` for virtual machine enhanced cloud disk, `CLOUD_BSSD` for virtual machine general purpose SSD cloud disk.
* `memory` - (Required, Int) Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
* `name` - (Required, String) Name of the SQL Server basic instance.
* `storage` - (Required, Int) Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
* `auto_renew` - (Optional, Int) Automatic renewal sign. 0 for normal renewal, 1 for automatic renewal, the default is 1 automatic renewal. Only valid when purchasing a prepaid instance.
* `auto_voucher` - (Optional, Int) Whether to use the voucher automatically; 1 for yes, 0 for no, the default is 0.
* `availability_zone` - (Optional, String, ForceNew) Availability zone.
* `charge_type` - (Optional, String, ForceNew) Pay type of the SQL Server basic instance. For now, only `POSTPAID_BY_HOUR` is valid.
* `collation` - (Optional, String) System character set sorting rule, default: Chinese_PRC_CI_AS.
* `engine_version` - (Optional, String, ForceNew) Version of the SQL Server basic database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
* `maintenance_start_time` - (Optional, String) Start time of the maintenance in one day, format like `HH:mm`.
* `maintenance_time_span` - (Optional, Int) The timespan of maintenance in one day, unit is hour.
* `maintenance_week_set` - (Optional, Set: [`Int`]) A list of integer indicates weekly maintenance. For example, [1,7] presents do weekly maintenance on every Monday and Sunday.
* `period` - (Optional, Int) Purchase instance period, the default value is 1, which means one month. The value does not exceed 48.
* `project_id` - (Optional, Int) Project ID, default value is 0.
* `security_groups` - (Optional, Set: [`String`]) Security group bound to the instance.
* `subnet_id` - (Optional, String, ForceNew) ID of subnet.
* `tags` - (Optional, Map) The tags of the SQL Server basic instance.
* `voucher_ids` - (Optional, Set: [`String`]) An array of voucher IDs, currently only one can be used for a single order.
* `vpc_id` - (Optional, String, ForceNew) ID of VPC.

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
$ terraform import tencentcloud_sqlserver_basic_instance.example mssql-3cdq7kx5
```

