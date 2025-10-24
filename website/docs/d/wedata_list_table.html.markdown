---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_list_table"
sidebar_current: "docs-tencentcloud-datasource-wedata_list_table"
description: |-
  Use this data source to query detailed information of WeData list table
---

# tencentcloud_wedata_list_table

Use this data source to query detailed information of WeData list table

## Example Usage

```hcl
data "tencentcloud_wedata_list_table" "example" {}
```

## Argument Reference

The following arguments are supported:

* `catalog_name` - (Optional, String) Directory name.
* `database_name` - (Optional, String) Database name.
* `datasource_id` - (Optional, Int) Data source ID.
* `keyword` - (Optional, String) Table search keyword.
* `result_output_file` - (Optional, String) Used to save results.
* `schema_name` - (Optional, String) Database schema name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Schema record list.
  * `business_metadata` - Business metadata of the table.
    * `tag_names` - Tag names.
  * `create_time` - Creation time.
  * `database_name` - Database name.
  * `description` - Data table description.
  * `guid` - Data table GUID.
  * `name` - Data table name.
  * `schema_name` - Database schema name.
  * `table_type` - Table type.
  * `technical_metadata` - Technical metadata of the table.
    * `location` - Data table location.
    * `owner` - Owner.
    * `storage_size` - Storage size.
  * `update_time` - Update time.


