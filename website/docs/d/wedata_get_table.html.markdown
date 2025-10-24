---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_get_table"
sidebar_current: "docs-tencentcloud-datasource-wedata_get_table"
description: |-
  Use this data source to query detailed information of WeData get table
---

# tencentcloud_wedata_get_table

Use this data source to query detailed information of WeData get table

## Example Usage

```hcl
data "tencentcloud_wedata_get_table" "example" {
  table_guid = "ktDR4ymhp2_nlfClXhwxRQ"
}
```

## Argument Reference

The following arguments are supported:

* `table_guid` - (Required, String) Table GUID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data table details.
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
    * `owner` - Responsible person.
    * `storage_size` - Storage size.
  * `update_time` - Update time.


