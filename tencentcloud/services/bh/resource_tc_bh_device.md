Provides a resource to create a BH device

Example Usage

```hcl
resource "tencentcloud_bh_device" "example" {
  device_set {
    os_name = "Linux"
    ip      = "1.1.1.1"
    port    = 22
    name    = "tf-example"
  }
}
```

Import

BH device can be imported using the id, e.g.

```
terraform import tencentcloud_bh_device.example 1875
```
