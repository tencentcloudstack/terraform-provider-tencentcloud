---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_cross_region_zone"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_cross_region_zone"
description: |-
  Use this data source to query detailed information of sqlserver datasource_cross_region_zone
---

# tencentcloud_sqlserver_cross_region_zone

Use this data source to query detailed information of sqlserver datasource_cross_region_zone

## Example Usage

```hcl
data "tencentcloud_sqlserver_cross_region_zone" "example" {
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of mssql-j8kv137v.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region` - The string ID of the region where the standby machine is located, such as: ap-guangzhou.
* `zone` - The string ID of the availability zone where the standby machine is located, such as: ap-guangzhou-1.


