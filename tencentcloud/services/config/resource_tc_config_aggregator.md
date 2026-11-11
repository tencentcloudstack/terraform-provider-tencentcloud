Provides a resource to create a Config aggregator.

Example Usage

```hcl
resource "tencentcloud_config_aggregator" "example" {
  name        = "tf-example-aggregator"
  description = "tf example aggregator"
  type        = "CUSTOM"
  owner_uin   = "100012345678"

  aggregator_accounts {
    member_uin  = 100012345679
    member_name = "member-1"
  }

  aggregator_accounts {
    member_uin  = 100012345680
    member_name = "member-2"
  }
}
```

Import

Config aggregator can be imported using the composite id, e.g. `account_group_id#owner_uin`.

```
terraform import tencentcloud_config_aggregator.example ca-xxxxxxxx#100012345678
```
