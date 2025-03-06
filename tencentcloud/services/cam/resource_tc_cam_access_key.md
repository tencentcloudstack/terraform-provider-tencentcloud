Provides a resource to create a CAM access key

Example Usage

Create access key

```hcl
data "tencentcloud_user_info" "info" {}

resource "tencentcloud_cam_access_key" "example" {
  target_uin = data.tencentcloud_user_info.info.uin
}
```

Update access key

```hcl
data "tencentcloud_user_info" "info" {}

resource "tencentcloud_cam_access_key" "example" {
  target_uin = data.tencentcloud_user_info.info.uin
  status     = "Inactive"
}
```

Encrypted access key

```hcl
data "tencentcloud_user_info" "info" {}

resource "tencentcloud_cam_access_key" "example" {
  target_uin = data.tencentcloud_user_info.info.uin
  pgp_key    = "keybase:some_person_that_exists"
}
```

Import

cam access key can be imported using the id, e.g.

```
terraform import tencentcloud_cam_access_key.example 100037718101#AKID7F******************
```