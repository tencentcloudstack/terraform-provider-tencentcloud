---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_get_table_columns"
sidebar_current: "docs-tencentcloud-datasource-wedata_get_table_columns"
description: |-
  Use this data source to query detailed information of WeData get table columns
---

# tencentcloud_wedata_get_table_columns

Use this data source to query detailed information of WeData get table columns

## Example Usage

```hcl
data "tencentcloud_wedata_get_table_columns" "example" {
  table_guid = "ktDR4ymhp2_nlfClXhwxRQ"
}
```

## Argument Reference

The following arguments are supported:

* `table_guid` - (Required, String) Table GUID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Table column list.
  * `description` - Field description.
  * `is_partition` - Whether it is a partition field.
  * `length` - Field length.
  * `name` - Field name.
  * `position` - Field position.
  * `type` - Field type.


