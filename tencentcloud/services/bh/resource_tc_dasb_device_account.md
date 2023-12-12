Provides a resource to create a dasb device_account

Example Usage

```hcl
resource "tencentcloud_dasb_device_account" "example" {
  device_id = 100
  account   = "root"
}
```

Import

dasb device_account can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_account.example 11
```