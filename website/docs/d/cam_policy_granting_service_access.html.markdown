---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_policy_granting_service_access"
sidebar_current: "docs-tencentcloud-datasource-cam_policy_granting_service_access"
description: |-
  Use this data source to query detailed information of cam policy_granting_service_access
---

# tencentcloud_cam_policy_granting_service_access

Use this data source to query detailed information of cam policy_granting_service_access

## Example Usage

```hcl
data "tencentcloud_cam_policy_granting_service_access" "policy_granting_service_access" {
  role_id      = 4611686018436805021
  service_type = "cam"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, Int) Group Id, one of the three (TargetUin, RoleId, GroupId) must be passed.
* `result_output_file` - (Optional, String) Used to save results.
* `role_id` - (Optional, Int) Role Id, one of the three (TargetUin, RoleId, GroupId) must be passed.
* `service_type` - (Optional, String) Service type, this field needs to be passed when viewing the details of the service authorization interface.
* `target_uin` - (Optional, Int) Sub-account uin, one of the three (TargetUin, RoleId, GroupId) must be passed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List.
  * `action` - Action list.
    * `description` - Action description.
    * `name` - Action name.
  * `policy` - Policy list.
    * `policy_description` - Policy description.
    * `policy_id` - Policy Id.
    * `policy_name` - Policy name.
    * `policy_type` - Polic type.
  * `service` - Service info.
    * `service_name` - Service name.
    * `service_type` - Service type.


