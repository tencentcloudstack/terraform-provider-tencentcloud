---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_mfa_flag"
sidebar_current: "docs-tencentcloud-resource-cam_mfa_flag"
description: |-
  Provides a resource to create a cam mfa_flag
---

# tencentcloud_cam_mfa_flag

Provides a resource to create a cam mfa_flag

## Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

resource "tencentcloud_cam_mfa_flag" "mfa_flag" {
  op_uin = data.tencentcloud_user_info.info.uin
  login_flag {
    phone  = 0
    stoken = 1
    wechat = 0
  }
  action_flag {
    phone  = 0
    stoken = 1
    wechat = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `op_uin` - (Required, Int) Operate uin.
* `action_flag` - (Optional, List) Action flag setting.
* `login_flag` - (Optional, List) Login flag setting.

The `action_flag` object supports the following:

* `phone` - (Optional, Int) Phone.
* `stoken` - (Optional, Int) Soft token.
* `wechat` - (Optional, Int) Wechat.

The `login_flag` object supports the following:

* `phone` - (Optional, Int) Phone.
* `stoken` - (Optional, Int) Soft token.
* `wechat` - (Optional, Int) Wechat.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cam mfa_flag can be imported using the id, e.g.

```
terraform import tencentcloud_cam_mfa_flag.mfa_flag mfa_flag_id
```

