Provides a resource to create a BH access white list config

Example Usage

```hcl
resource "tencentcloud_bh_access_white_list_config" "example" {
  allow_any  = false
  allow_auto = false
}
```

Import

BH access white list config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_bh_access_white_list_config.example zDxkr768TFYadnFdX1fusQ==
```
