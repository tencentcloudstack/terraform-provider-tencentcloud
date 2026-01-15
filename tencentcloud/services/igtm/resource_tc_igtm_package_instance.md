Provides a resource to create a IGTM package instance

Example Usage

```hcl
resource "tencentcloud_igtm_package_instance" "example" {
  goods_type   = "STANDARD"
  auto_renew   = 1
  time_span    = 1
  auto_voucher = 1
}
```

Import

IGTM package instance can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_package_instance.example ins-wtqicjwzzze
```
