Provides a resource to create a cam access_key

Example Usage

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
}
```
Update
```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  status = "Inactive"
}
```
Encrypted
```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  pgp_key = "keybase:some_person_that_exists"
}
```
Import

cam access_key can be imported using the id, e.g.

```
terraform import tencentcloud_cam_access_key.access_key access_key_id
```