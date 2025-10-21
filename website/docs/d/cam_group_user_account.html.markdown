---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_group_user_account"
sidebar_current: "docs-tencentcloud-datasource-cam_group_user_account"
description: |-
  Use this data source to query detailed information of cam group_user_account
---

# tencentcloud_cam_group_user_account

Use this data source to query detailed information of cam group_user_account

## Example Usage

```hcl
data "tencentcloud_cam_group_user_account" "group_user_account" {
  sub_uin = 100033690181
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `rp` - (Optional, Int) Number per page. The default is 20.
* `sub_uin` - (Optional, Int) Sub-user uin.
* `uid` - (Optional, Int) Sub-user uid.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_info` - User group information.
  * `create_time` - Create time.
  * `group_id` - User group ID.
  * `group_name` - User group name.
  * `remark` - Remark.
* `total_num` - The total number of user groups the sub-user has joined.


