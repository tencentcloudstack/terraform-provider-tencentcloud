Provides a resource to create a tse cngw_strategy_bind_group

Example Usage

```hcl
resource "tencentcloud_tse_cngw_strategy_bind_group" "cngw_strategy_bind_group" {
  gateway_id  = "gateway-cf8c99c3"
  strategy_id = "strategy-806ea0dd"
  group_id    = "group-a160d123"
  option      = "bind"
}
```

Import

tse cngw_strategy_bind_group can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group cngw_strategy_bind_group_id
```
