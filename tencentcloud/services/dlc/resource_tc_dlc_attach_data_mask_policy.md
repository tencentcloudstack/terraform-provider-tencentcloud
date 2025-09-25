Provides a resource to create a DLC dlc attach data mask policy

Example Usage

```hcl
resource "tencentcloud_dlc_attach_data_mask_policy" "example" {
  data_mask_strategy_policy_set {
    policy_info {
      database    = ""
      catalog     = ""
      table       = ""
      operation   = ""
      policy_type = ""
      column      = ""
      mode        = ""
    }

    data_mask_strategy_id = ""
    column_type           = ""
  }
}
```
