Provides a resource to create a cam mfa_flag

Example Usage

```hcl
data "tencentcloud_user_info" "info"{}

resource "tencentcloud_cam_mfa_flag" "mfa_flag" {
  op_uin = data.tencentcloud_user_info.info.uin
  login_flag {
	phone = 0
	stoken = 1
	wechat = 0
  }
  action_flag {
	phone = 0
	stoken = 1
	wechat = 0
  }
}

```

Import

cam mfa_flag can be imported using the id, e.g.

```
terraform import tencentcloud_cam_mfa_flag.mfa_flag mfa_flag_id
```