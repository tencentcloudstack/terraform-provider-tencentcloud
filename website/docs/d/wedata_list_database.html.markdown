---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_list_database"
sidebar_current: "docs-tencentcloud-datasource-wedata_list_database"
description: |-
  Use this data source to query detailed information of WeData list database
---

# tencentcloud_wedata_list_database

Use this data source to query detailed information of WeData list database

## Example Usage

```hcl
data "tencentcloud_wedata_list_database" "example" {}
```

## Argument Reference

The following arguments are supported:

* `catalog_name` - (Optional, String) Catalog name.
* `datasource_id` - (Optional, Int) Data source ID.
* `keyword` - (Optional, String) Database name search keyword.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Database record list.
  * `catalog_name` - Database catalog.
  * `description` - Database description.
  * `guid` - Database GUID.
  * `location` - Database location.
  * `name` - Database name.
  * `storage_size` - Database storage size.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [ListDatabase](https://cloud.tencent.com/document/api/1267/123992) - This API is used to obtain asset database info.

