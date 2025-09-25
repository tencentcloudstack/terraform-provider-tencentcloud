Provides a resource to create a DLC data mask strategy

Example Usage

```hcl
resource "tencentcloud_dlc_data_mask_strategy" "example" {
  strategy {
    strategy_name = "tf-example"
    strategy_type = "tf-example"
    strategy_desc = "description."
    groups = {
      work_group_id = 221498
      strategy_type = "MASK"
    }
    users = ""
  }
}
```

Import

DLC data mask strategy can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_data_mask_strategy.example strategyId
```
