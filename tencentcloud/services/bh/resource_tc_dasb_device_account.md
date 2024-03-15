Provides a resource to create a dasb device_account

Example Usage

```hcl
resource "tencentcloud_dasb_device" "example" {
  os_name       = "Linux"
  ip            = "192.168.0.1"
  port          = 80
  name          = "tf_example"
}

resource "tencentcloud_dasb_device_account" "example" {
  device_id = tencentcloud_dasb_device.example.id
  account   = "root"
}
```

Import

dasb device_account can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_account.example 11
```