---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_list_schema"
sidebar_current: "docs-tencentcloud-datasource-wedata_list_schema"
description: |-
  Use this data source to query detailed information of WeData list schema
---

# tencentcloud_wedata_list_schema

Use this data source to query detailed information of WeData list schema

## Example Usage

```hcl
data "tencentcloud_wedata_list_schema" "example" {}
```

## Argument Reference

The following arguments are supported:

* `catalog_name` - (Optional, String) Catalog name.
* `database_name` - (Optional, String) Database name.
* `datasource_id` - (Optional, Int) Data source ID.
* `keyword` - (Optional, String) Database schema search keyword.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Schema record list.
  * `database_name` - Database name.
  * `guid` - Schema GUID.
  * `name` - Schema name.


