---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_collation_time_zone"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_collation_time_zone"
description: |-
  Use this data source to query detailed information of sqlserver datasource_collation_time_zone
---

# tencentcloud_sqlserver_collation_time_zone

Use this data source to query detailed information of sqlserver datasource_collation_time_zone

## Example Usage

```hcl
data "tencentcloud_sqlserver_collation_time_zone" "sqlserver_collation_time_zone" {
  machine_type = "CLOUD_PREMIUM"
}
```

## Argument Reference

The following arguments are supported:

* `db_version` - (Optional, String) Purchase instance version number.
* `machine_type` - (Optional, String) The host type of the purchased instance, PM-physical machine, CLOUD_PREMIUM-cloud server high-performance cloud disk, CLOUD_SSD-cloud server SSD cloud disk, CLOUD_HSSD-cloud server enhanced SSD cloud disk, CLOUD_TSSD-cloud server ultra-fast SSD cloud disk, CLOUD_BSSD-cloud server general SSD cloud disk, CLOUD_BASIC-cloud server cloud hard disk, the default value is PM.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `collation` - List of system character set collations.
* `time_zone` - System time zone list.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeCollationTimeZone](https://www.tencentcloud.com/document/api/238/72700)

