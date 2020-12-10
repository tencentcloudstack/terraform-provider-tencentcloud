---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_readonly_groups"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_readonly_groups"
description: |-
  Use this data source to query the list of SQL Server readonly groups.
---

# tencentcloud_sqlserver_readonly_groups

Use this data source to query the list of SQL Server readonly groups.

## Example Usage

```hcl
data "tencentcloud_sqlserver_readonly_groups" "master" {
  master_instance_id = "mssql-3cdq7kx5"
}
```

## Argument Reference

The following arguments are supported:

* `master_instance_id` - (Optional) Master SQL Server instance ID.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of SQL Server readonly group. Each element contains the following attributes:
  * `id` - ID of the readonly group.
  * `is_offline_delay` - Indicate whether to offline delayed readonly instances.
  * `master_instance_id` - Master instance id.
  * `max_delay_time` - Maximum delay time of the readonly instances.
  * `min_instances` - Minimum readonly instances that stays in the group.
  * `name` - Name of the readonly group.
  * `readonly_instance_set` - Readonly instance ID set of the readonly group.
  * `status` - Status of the readonly group. `1` for running, `5` for applying.
  * `vip` - Virtual IP address of the readonly group.
  * `vport` - Virtual port of the readonly group.


