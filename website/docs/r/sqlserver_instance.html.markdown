---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_instance"
description: |-
  Use this resource to create SQL Server instance
---

# tencentcloud_sqlserver_instance

Use this resource to create SQL Server instance

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

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 16
  storage           = 100
}
```

## Argument Reference

The following arguments are supported:

* `memory` - (Required, Int) Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
* `name` - (Required, String) Name of the SQL Server instance.
* `storage` - (Required, Int) Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
* `auto_renew` - (Optional, Int) Automatic renewal sign. 0 for normal renewal, 1 for automatic renewal (Default). Only valid when purchasing a prepaid instance.
* `auto_voucher` - (Optional, Int) Whether to use the voucher automatically; 1 for yes, 0 for no, the default is 0.
* `availability_zone` - (Optional, String, ForceNew) Availability zone.
* `charge_type` - (Optional, String, ForceNew) Pay type of the SQL Server instance. Available values `PREPAID`, `POSTPAID_BY_HOUR`.
* `engine_version` - (Optional, String, ForceNew) Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
* `ha_type` - (Optional, String, ForceNew, **Deprecated**) It has been deprecated from version 1.81.136. Instance type. `DUAL` (dual-server high availability), `CLUSTER` (cluster). Default is `DUAL`.
* `maintenance_start_time` - (Optional, String) Start time of the maintenance in one day, format like `HH:mm`.
* `maintenance_time_span` - (Optional, Int) The timespan of maintenance in one day, unit is hour.
* `maintenance_week_set` - (Optional, Set: [`Int`]) A list of integer indicates weekly maintenance. For example, [2,7] presents do weekly maintenance on every Tuesday and Sunday.
* `multi_zones` - (Optional, Bool, ForceNew) Indicate whether to deploy across availability zones.
* `period` - (Optional, Int) Purchase instance period in month. The value does not exceed 48.
* `project_id` - (Optional, Int) Project ID, default value is 0.
* `security_groups` - (Optional, Set: [`String`]) Security group bound to the instance.
* `subnet_id` - (Optional, String) ID of subnet.
* `tags` - (Optional, Map) The tags of the SQL Server.
* `time_zone` - (Optional, String) System time zone, default: `China Standard Time`.
* `voucher_ids` - (Optional, Set: [`String`]) An array of voucher IDs, currently only one can be used for a single order.
* `vpc_id` - (Optional, String) ID of VPC.
* `wait_switch` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.81.2. The way to execute the allocation. Supported values include: 0 - execute immediately, 1 - execute in maintenance window.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server instance.
* `dns_pod_domain` - Internet address domain name.
* `ro_flag` - Readonly flag. `RO` (read-only instance), `MASTER` (primary instance with read-only instances). If it is left empty, it refers to an instance which is not read-only and has no RO group.
* `status` - Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
* `tgw_wan_vport` - External port number.
* `vip` - IP for private access.
* `vport` - Port for private access.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDBInstances](https://cloud.tencent.com/document/api/238/19973) - This API is used to create high-availability instances (local disk).
* Create: [DescribeFlowStatus](https://cloud.tencent.com/document/api/238/19967) - This API is used to query flow status.
* Create: [DescribeOrders](https://cloud.tencent.com/document/api/238/19937) - This API is used to query order information.
* Create: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/238/19958) - This API is used to rename an instance.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeMaintenanceSpan](https://cloud.tencent.com/document/api/238/46283) - This API is used to query the maintenance time window of an instance based on its instance ID.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeFlowStatus](https://cloud.tencent.com/document/api/238/19967) - This API is used to query flow status.
* Update: [ModifyDBInstanceNetwork](https://cloud.tencent.com/document/api/238/52347) - This API is used to switch a running instance from a VPC to another.
* Update: [ModifyDBInstanceProject](https://cloud.tencent.com/document/api/238/19957) - This API is used to modify the project to which a database instance belongs.
* Update: [ModifyMaintenanceSpan](https://cloud.tencent.com/document/api/238/46281) - This API is used to modify the maintenance window of the instance.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteDBInstance](https://cloud.tencent.com/document/api/238/46806) - This API is used to release SQL server instances (eliminated immediately) in the recycle bin. The released instances will be physically terminated after being retained for a period of time. Their publish-subscribe relationships will be automatically disassociated, and their RO replicas will be automatically released.
* Delete: [DescribeDBInstances](https://cloud.tencent.com/document/api/238/19969) - This API is used to query the list of instances.
* Delete: [TerminateDBInstance](https://cloud.tencent.com/document/api/238/38210) - This API is used to isolate an instance to move it into a recycle bin.

## Import

SQL Server instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_instance.example mssql-3cdq7kx5
```

