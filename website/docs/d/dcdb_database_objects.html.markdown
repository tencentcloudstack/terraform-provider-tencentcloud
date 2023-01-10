---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_database_objects"
sidebar_current: "docs-tencentcloud-datasource-dcdb_database_objects"
description: |-
  Use this data source to query detailed information of dcdb database_objects
---

# tencentcloud_dcdb_database_objects

Use this data source to query detailed information of dcdb database_objects

## Example Usage

```hcl
data "tencentcloud_dcdb_database_objects" "database_objects" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name     = & lt ; nil & gt ;
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) Database name, obtained through the DescribeDatabases api.
* `instance_id` - (Required, String) The ID of instance.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `funcs` - Function list.
  * `func` - The name of function.
* `procs` - Procedure list.
  * `proc` - The name of procedure.
* `tables` - Table list.
  * `table` - The name of table.
* `views` - View list.
  * `view` - The name of view.


