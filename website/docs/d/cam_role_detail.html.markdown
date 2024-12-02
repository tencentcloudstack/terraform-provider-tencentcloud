---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_detail"
sidebar_current: "docs-tencentcloud-datasource-cam_role_detail"
description: |-
  Use this data source to query detailed information of cam role detail
---

# tencentcloud_cam_role_detail

Use this data source to query detailed information of cam role detail

## Example Usage

### Query cam role detail by role ID

```hcl
data "tencentcloud_cam_role_detail" "example" {
  role_id = "4611686018441060141"
}
```

### Query cam role detail by role name

```hcl
data "tencentcloud_cam_role_detail" "example" {
  role_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `role_id` - (Optional, String) Role ID, used to specify role. Input either `RoleId` or `RoleName`.
* `role_name` - (Optional, String) Role name, used to specify role. Input either `RoleId` or `RoleName`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_info` - Role details.


