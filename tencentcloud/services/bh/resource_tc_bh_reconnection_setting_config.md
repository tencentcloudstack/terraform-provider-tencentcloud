Provides a resource to create a BH reconnection setting config

Example Usage

```hcl
resource "tencentcloud_bh_reconnection_setting_config" "example" {
  reconnection_max_count = 5
  enable                 = true
}
```

Import

BH reconnection setting config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_bh_reconnection_setting_config.example gO1Ew6OEgLcQun164XiWmw==
```
