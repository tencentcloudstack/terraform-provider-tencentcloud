Use this data source to query detailed information of CAM users

Example Usage

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