---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_top_space_table_time_series"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_top_space_table_time_series"
description: |-
  Use this data source to query detailed information of dbbrain top_space_table_time_series
---

# tencentcloud_dbbrain_top_space_table_time_series

Use this data source to query detailed information of dbbrain top_space_table_time_series

## Example Usage

```hcl
data "tencentcloud_dbbrain_top_space_table_time_series" "top_space_table_time_series" {
  instance_id = "%s"
  sort_by     = "DataLength"
  start_date  = "%s"
  end_date    = "%s"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `end_date` - (Optional, String) The deadline, such as 2021-01-01, the earliest is the 29th day before the current day, and the default is the current day.
* `limit` - (Optional, Int) The number of Top tables returned, the maximum value is 100, and the default is 20.
* `product` - (Optional, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) The sorting field used to filter the Top table. The optional fields include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, and PhysicalFileSize. The default is PhysicalFileSize.
* `start_date` - (Optional, String) The start date, such as 2021-01-01, the earliest is the 29th day before the current day, and the default is the 6th day before the deadline.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `top_space_table_time_series` - The time-series data list of the returned Top tablespace statistics.
  * `engine` - Storage engine for database tables.
  * `series_data` - Spatial index data in unit time interval.
    * `series` - Monitor metrics.
      * `metric` - Indicator name.
      * `unit` - Indicator unit.
      * `values` - Index value. Note: This field may return null, indicating that no valid value can be obtained.
    * `timestamp` - The timestamp corresponding to the monitoring indicator.
  * `table_name` - table name.
  * `table_schema` - databases name.


