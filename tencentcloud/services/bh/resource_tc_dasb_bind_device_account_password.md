Provides a resource to create a dasb bind_device_account_password

Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_account_password" "example" {
  device_account_id = 16
  password          = "TerraformPassword"
}
```