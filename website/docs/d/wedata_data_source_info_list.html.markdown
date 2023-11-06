---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_data_source_info_list"
sidebar_current: "docs-tencentcloud-datasource-wedata_data_source_info_list"
description: |-
  Use this data source to query detailed information of wedata data_source_info_list
---

# tencentcloud_wedata_data_source_info_list

Use this data source to query detailed information of wedata data_source_info_list

## Example Usage

```hcl
data "tencentcloud_wedata_data_source_info_list" "example" {
  project_id = "1927766435649077248"
  filters {
    name   = "Name"
    values = ["tf_example"]
  }

  order_fields {
    name      = "CreateTime"
    direction = "DESC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `datasource_name` - (Optional, String) DatasourceName.
* `filters` - (Optional, List) Filters.
* `order_fields` - (Optional, List) OrderFields.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, String) Datasource type.

The `filters` object supports the following:

* `name` - (Optional, String) Filter name.
* `values` - (Optional, Set) Filter values.

The `order_fields` object supports the following:

* `direction` - (Required, String) OrderFields rule.
* `name` - (Required, String) OrderFields name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `datasource_set` - DatasourceSet.
  * `cluster_id` - ClusterId.
  * `database_names` - DatabaseNames.
  * `description` - Description.
  * `id` - Id.
  * `instance` - Instance.
  * `name` - Name.
  * `region` - Region.
  * `type` - Type.
  * `version` - Desc.


