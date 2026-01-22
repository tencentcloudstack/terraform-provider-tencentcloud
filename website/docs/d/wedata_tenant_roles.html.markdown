---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_tenant_roles"
sidebar_current: "docs-tencentcloud-datasource-wedata_tenant_roles"
description: |-
  Use this data source to query detailed information of WeData tenant roles
---

# tencentcloud_wedata_tenant_roles

Use this data source to query detailed information of WeData tenant roles

## Example Usage

```hcl
data "tencentcloud_wedata_tenant_roles" "example" {
  role_display_name = "tf_example"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `role_display_name` - (Optional, String) Role Chinese display name fuzzy search, can only pass one value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Main account role list.
  * `description` - Description.
  * `role_display_name` - Role display name.
  * `role_id` - Role ID.
  * `role_name` - Role name.


