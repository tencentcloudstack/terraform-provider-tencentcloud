---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_code_max_permission"
sidebar_current: "docs-tencentcloud-datasource-wedata_code_max_permission"
description: |-
  Use this data source to query detailed information of WeData code max permission
---

# tencentcloud_wedata_code_max_permission

Use this data source to query detailed information of WeData code max permission

## Example Usage

```hcl
data "tencentcloud_wedata_code_max_permission" "example" {
  project_id  = "3108707295180644352"
  resource_id = "f0c14b9d-003e-4325-8830-d1a9fa934ed6"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `resource_id` - (Required, String) Unique ID of authorization resource, folder ID or file ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - User's recursive maximum permission type for CodeStudio files/folders.
  * `permission_type` - Authorization permission type (CAN_VIEW/CAN_RUN/CAN_EDIT/CAN_MANAGE).


