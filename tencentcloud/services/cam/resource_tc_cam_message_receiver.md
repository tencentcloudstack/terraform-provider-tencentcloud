Provides a resource to create a CAM message receiver

~> **NOTE:** For security reasons, the CAM will return the `email` and `phone_number` parameter values in encrypted form. Please use the `ignore_changes` function in Terraform's `lifecycle` to include these two parameters.

Example Usage

```hcl
resource "tencentcloud_cam_message_receiver" "example" {
  name         = "tf-example"
  remark       = "remark."
  country_code = "86"
  phone_number = "18123456789"
  email        = "demo@qq.com"

  lifecycle {
    ignore_changes = [ email, phone_number ]
  }
}
```

Import

CAM message receiver can be imported using the id, e.g.

```
terraform import tencentcloud_cam_message_receiver.example tf-example
```
