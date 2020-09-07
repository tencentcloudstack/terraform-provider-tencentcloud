---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_zone_config"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_zone_config"
description: |-
  Use this data source to query purchasable specification configuration for each availability zone in this specific region.
---

# tencentcloud_sqlserver_zone_config

Use this data source to query purchasable specification configuration for each availability zone in this specific region.

## Example Usage

```hcl
data "tencentcloud_sqlserver_zone_config" "mysqlserver" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_list` - A list of availability zones. Each element contains the following attributes:
  * `availability_zone` - Alphabet ID of availability zone.
  * `specinfo_list` - A list of specinfo configurations for the specific availability zone. Each element contains the following attributes:
    * `charge_type` - Billing mode under this specification. Valid values are `POSTPAID_BY_HOUR`, `PREPAID` and `ALL` which means both POSTPAID_BY_HOUR and PREPAID.
    * `cpu` - Number of CPU cores.
    * `db_version_name` - Version name corresponding to the `db_version` field.
    * `db_version` - Database version information. Valid values: `2008R2 (SQL Server 2008 Enterprise)`, `2012SP3 (SQL Server 2012 Enterprise)`, `2016SP1 (SQL Server 2016 Enterprise)`, `201602 (SQL Server 2016 Standard)`, `2017 (SQL Server 2017 Enterprise)`.
    * `machine_type` - Model ID.
    * `max_storage_size` - Maximum disk size under this specification in GB.
    * `memory` - Memory size in GB.
    * `min_storage_size` - Minimum disk size under this specification in GB.
    * `qps` - QPS of this specification.
    * `spec_id` - Instance specification ID.
  * `zone_id` - Number ID of availability zone.


