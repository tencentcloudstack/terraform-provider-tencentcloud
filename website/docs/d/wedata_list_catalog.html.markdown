---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_list_catalog"
sidebar_current: "docs-tencentcloud-datasource-wedata_list_catalog"
description: |-
  Use this data source to query detailed information of WeData list catalog
---

# tencentcloud_wedata_list_catalog

Use this data source to query detailed information of WeData list catalog

## Example Usage

```hcl
data "tencentcloud_wedata_list_catalog" "example" {}
```

## Argument Reference

The following arguments are supported:

* `parent_catalog_id` - (Optional, String) Parent catalog ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Catalog record list.
  * `name` - Catalog name.


