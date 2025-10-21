---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_users"
sidebar_current: "docs-tencentcloud-datasource-cam_users"
description: |-
  Use this data source to query detailed information of CAM users
---

# tencentcloud_cam_users

Use this data source to query detailed information of CAM users

## Example Usage

```hcl
# query by name
data "tencentcloud_cam_users" "foo" {
  name = "cam-user-test"
}

# query by email
data "tencentcloud_cam_users" "bar" {
  email = "hello@test.com"
}

# query by phone
data "tencentcloud_cam_users" "far" {
  phone_num = "12345678910"
}
```

## Argument Reference

The following arguments are supported:

* `console_login` - (Optional, Bool) Indicate whether the user can login in.
* `country_code` - (Optional, String) Country code of the CAM user to be queried.
* `email` - (Optional, String) Email of the CAM user to be queried.
* `name` - (Optional, String) Name of CAM user to be queried.
* `phone_num` - (Optional, String) Phone num of the CAM user to be queried.
* `remark` - (Optional, String) Remark of the CAM user to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `uid` - (Optional, Int) Uid of the CAM user to be queried.
* `uin` - (Optional, Int) Uin of the CAM user to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_list` - A list of CAM users. Each element contains the following attributes:
  * `country_code` - Country code of the CAM user.
  * `email` - Email of the CAM user.
  * `name` - Name of CAM user.
  * `phone_num` - Phone num of the CAM user.
  * `remark` - Remark of the CAM user.
  * `uid` - Uid of the CAM user.
  * `uin` - Uin of the CAM user.
  * `user_id` - ID of CAM user. Its value equals to `name` argument.


