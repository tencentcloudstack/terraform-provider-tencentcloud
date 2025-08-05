Provides a resource to create a CAM message receiver

Example Usage

```hcl
resource "tencentcloud_cam_message_receiver" "example" {
  name         = "tf-example"
  remark       = "remark."
  country_code = "86"
  phone_number = "18123456789"
  email        = "demo@qq.com"
}
```

Import

CAM message receiver can be imported using the id, e.g.

```
terraform import tencentcloud_cam_message_receiver.example tf-example
```
