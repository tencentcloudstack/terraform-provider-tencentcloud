Provides a resource to create a DLC data mask strategy

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
```

Import

DLC data mask strategy can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_data_mask_strategy.example 2fcab650-11a8-44ef-bf58-19c22af601b6
```
