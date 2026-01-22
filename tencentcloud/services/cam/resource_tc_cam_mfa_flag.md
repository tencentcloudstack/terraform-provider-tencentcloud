Provides a resource to create a CAM mfa flag

Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

resource "tencentcloud_cam_mfa_flag" "example" {
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

Import

CAM mfa flag can be imported using the id, e.g.

```
terraform import tencentcloud_cam_mfa_flag.example 100037718110
```