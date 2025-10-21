---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_top_space_schemas"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_top_space_schemas"
description: |-
  Use this data source to query detailed information of dbbrain top_space_schemas
---

# tencentcloud_dbbrain_top_space_schemas

Use this data source to query detailed information of dbbrain top_space_schemas

## Example Usage

```hcl
data "tencentcloud_dbbrain_top_space_schemas" "top_space_schemas" {
  instance_id = "%s"
  sort_by     = "DataLength"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `limit` - (Optional, Int) The number of Top libraries to return, the maximum value is 100, and the default is 20.
* `product` - (Optional, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) The sorting field used to filter the Top library. The optional fields include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, and PhysicalFileSize (only supported by ApsaraDB for MySQL instances). The default for ApsaraDB for MySQL instances is PhysicalFileSize, and the default for other product instances is TotalLength.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `timestamp` - Timestamp (in seconds) when library space data is collected.
* `top_space_schemas` - The returned list of top library space statistics.
  * `data_free` - Fragmentation space (MB).
  * `data_length` - data space (MB).
  * `frag_ratio` - Fragmentation rate (%).
  * `index_length` - Index space (MB).
  * `physical_file_size` - The sum (MB) of the independent physical file sizes corresponding to all tables in the library. Note: This field may return null, indicating that no valid value can be obtained.
  * `table_rows` - Number of lines.
  * `table_schema` - library name.
  * `total_length` - Total space used (MB).


