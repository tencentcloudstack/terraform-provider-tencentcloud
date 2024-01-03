---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_meta_table"
sidebar_current: "docs-tencentcloud-datasource-oceanus_meta_table"
description: |-
  Use this data source to query detailed information of oceanus meta_table
---

# tencentcloud_oceanus_meta_table

Use this data source to query detailed information of oceanus meta_table

## Example Usage

```hcl
data "tencentcloud_oceanus_meta_table" "example" {
  work_space_id = "space-6w8eab6f"
  catalog       = "_dc"
  database      = "_db"
  table         = "tf_table"
}
```

## Argument Reference

The following arguments are supported:

* `catalog` - (Required, String) Catalog name.
* `database` - (Required, String) Database name.
* `table` - (Required, String) Table name.
* `work_space_id` - (Required, String) Unique identifier of the space.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Scene time.
* `ddl` - Table creation statement, encoded in Base64.For example,Q1JFQVRFIFRBQkxFIGRhdGFnZW5fc291cmNlX3RhYmxlICggCiAgICBpZCBJTlQsIAogICAgbmFtZSBTVFJJTkcgCikgV0lUSCAoCidjb25uZWN0b3InPSdkYXRhZ2VuJywKJ3Jvd3MtcGVyLXNlY29uZCcgPSAnMScKKTs=.
* `serial_id` - Unique identifier of the metadata table.


