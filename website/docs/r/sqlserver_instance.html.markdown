---
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
resource "tencentcloud_sqlserver_instance" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-409mvdvv"
  subnet_id         = "subnet-nf9n81ps"
  project_id        = 123
  memory            = 2
  storage           = 100
}
```

## Argument Reference

The following arguments are supported:

* `memory` - (Required) Memory size (in GB). Allowed value must be larger than `memory` that data source `tencentcloud_sqlserver_specinfos` provides.
* `name` - (Required) Name of the SQL Server instance.
* `storage` - (Required) Disk size (in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_sqlserver_specinfos` provides.
* `availability_zone` - (Optional, ForceNew) Availability zone.
* `charge_type` - (Optional, ForceNew) Pay type of the SQL Server instance. For now, only `POSTPAID_BY_HOUR` is valid.
* `engine_version` - (Optional, ForceNew) Version of the SQL Server database engine. Allowed values are `2008R2`(SQL Server 2008 Enerprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) and `2017`(SQL Server 2017 Enterprise). Default is `2008R2`.
* `project_id` - (Optional) Project ID, default value is 0.
* `subnet_id` - (Optional, ForceNew) ID of subnet.
* `vpc_id` - (Optional, ForceNew) ID of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server instance.
* `private_access_ip` - IP for private access.
* `private_access_port` - Port for private access.
* `status` - Status of the SQL Server instance. 1 for applying, 2 for running, 3 for running with limit, 4 for isolated, 5 for recycling, 6 for recycled, 7 for running with task, 8 for off-line, 9 for expanding, 10 for migrating, 11 for readonly, 12 for rebooting.
* `used_storage` - Used storage.


## Import

SQL Server instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_instance.foo mssql-3cdq7kx5
```

