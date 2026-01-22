---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_project_roles"
sidebar_current: "docs-tencentcloud-datasource-wedata_project_roles"
description: |-
  Use this data source to query detailed information of WeData project roles
---

# tencentcloud_wedata_project_roles

Use this data source to query detailed information of WeData project roles

## Example Usage

```hcl
data "tencentcloud_wedata_project_roles" "example" {
  project_id        = "2982667120655491072"
  role_display_name = "tf_example"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `role_display_name` - (Optional, String) Role Chinese display name fuzzy search, can only pass one value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Role information.
  * `description` - Description.
  * `role_display_name` - Role display name.
  * `role_id` - Role ID.
  * `role_name` - Role name.


