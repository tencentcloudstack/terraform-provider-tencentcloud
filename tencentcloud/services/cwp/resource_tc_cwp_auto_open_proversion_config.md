Provides a resource to create a CWP auto open proversion config

Example Usage

```hcl
resource "tencentcloud_cwp_auto_open_proversion_config" "example" {
  status                       = "OPEN"
  protect_type                 = "FLAGSHIP_PREPAY"
  auto_repurchase_renew_switch = 1
  auto_repurchase_switch       = 1
  repurchase_renew_switch      = 0
  auto_bind_rasp_switch        = 1
  auto_open_rasp_switch        = 1
  auto_downgrade_switch        = 0
}
```

Import

CWP auto open proversion config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cwp_auto_open_proversion_config.example 9n0PnuKR3sEmY8COuKVb7g==
```
