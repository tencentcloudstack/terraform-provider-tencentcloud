Provides a resource to create a DLC attach data mask policy

Example Usage

```hcl
resource "tencentcloud_dlc_data_mask_strategy" "example" {
  strategy {
    strategy_name = "tf-example"
    strategy_desc = "description."
    groups {
      work_group_id = 70220
      strategy_type = "MASK"
    }
  }
}

resource "tencentcloud_dlc_attach_data_mask_policy" "example" {
  data_mask_strategy_policy_set {
    policy_info {
      database    = "tf-example"
      catalog     = "DataLakeCatalog"
      table       = "tf-example"
      column      = "id"
    }

    data_mask_strategy_id = tencentcloud_dlc_data_mask_strategy.example.id
    column_type           = "string"
  }
}
```
