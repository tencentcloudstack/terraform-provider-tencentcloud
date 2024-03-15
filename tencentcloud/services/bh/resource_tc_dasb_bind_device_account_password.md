Provides a resource to create a dasb bind_device_account_password

Example Usage

```hcl
resource "tencentcloud_dasb_device" "example" {
  os_name = "Linux"
  ip      = "192.168.0.1"
  port    = 80
  name    = "tf_example"
}

resource "tencentcloud_dasb_device_account" "example" {
  device_id = tencentcloud_dasb_device.example.id
  account   = "root"
}

resource "tencentcloud_dasb_bind_device_account_password" "example" {
  device_account_id = tencentcloud_dasb_device_account.example.id
  password          = "TerraformPassword"
}
```