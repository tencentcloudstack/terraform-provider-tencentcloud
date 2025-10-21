---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_database_tables"
sidebar_current: "docs-tencentcloud-datasource-dcdb_database_tables"
description: |-
  Use this data source to query detailed information of dcdb database_tables
---

# tencentcloud_dcdb_database_tables

Use this data source to query detailed information of dcdb database_tables

## Example Usage

```hcl
data "tencentcloud_dcdb_database_tables" "database_tables" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name     = & lt ; nil & gt ;
  table       = & lt ; nil & gt ;
  table       = & lt ; nil & gt ;
  cols {
    col  = & lt ; nil & gt ;
    type = & lt ; nil & gt ;

  }
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) Database name, obtained through the DescribeDatabases api.
* `instance_id` - (Required, String) The ID of instance.
* `table` - (Required, String) Table name, obtained through the DescribeDatabaseObjects api.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cols` - Column information.
  * `col` - The name of column.
  * `type` - Column type.


