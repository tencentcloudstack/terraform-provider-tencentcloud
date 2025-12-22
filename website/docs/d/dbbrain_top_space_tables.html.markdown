---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_top_space_tables"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_top_space_tables"
description: |-
  Use this data source to query detailed information of dbbrain top_space_tables
---

# tencentcloud_dbbrain_top_space_tables

Use this data source to query detailed information of dbbrain top_space_tables

## Example Usage

Sort by PhysicalFileSize

```hcl
data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by     = "PhysicalFileSize"
  product     = "mysql"
}
```

Sort by TotalLength

```hcl
data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by     = "PhysicalFileSize"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `limit` - (Optional, Int) The number of Top tables returned, the maximum value is 100, and the default is 20.
* `product` - (Optional, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) The sorting field used to filter the Top table. The optional fields include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, and PhysicalFileSize (only supported by ApsaraDB for MySQL instances). The default for ApsaraDB for MySQL instances is PhysicalFileSize, and the default for other product instances is TotalLength.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `timestamp` - The timestamp (in seconds) of collecting tablespace data.
* `top_space_tables` - The list of Top tablespace statistics returned.
  * `data_free` - Fragmentation space (MB).
  * `data_length` - data space (MB).
  * `engine` - Storage engine for database tables.
  * `frag_ratio` - Fragmentation rate (%).
  * `index_length` - Index space (MB).
  * `physical_file_size` - The independent physical file size (MB) corresponding to the table.
  * `table_name` - table name.
  * `table_rows` - Number of lines.
  * `table_schema` - database name.
  * `total_length` - Total space used (MB).


