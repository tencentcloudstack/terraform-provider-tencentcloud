Provides a resource to create a dasb device

Example Usage

```hcl
resource "tencentcloud_dasb_device" "example" {
  os_name       = "Linux"
  ip            = "192.168.0.1"
  port          = 80
  name          = "tf_example"
  department_id = "1.2.3"
}
```

Import

dasb device can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device.example 17
```