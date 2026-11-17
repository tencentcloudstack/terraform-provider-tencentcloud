Provides a resource to create a FWM (Firewall Manager) strategy.

Example Usage

Create pre strategy

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product = "cfw_edge_acl"
  receive_account = [
    "100011949846",
  ]
  receive_group = [
    "acg-y5tocz3d",
    "acg-e6cvbls1"
  ]

  pre_strategy {
    group_id = "fwmrg-qww1luux7z"
  }
}
```

Create post strategy

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product = "cfw_edge_acl"
  receive_account = [
    "100011949846",
    "100049315301"
  ]
  receive_group = [
    "acg-y5tocz3d",
  ]

  post_strategy {
    group_id = "fwmrg-qww1luux7z"
  }
}
```

Import

FWM strategy can be imported using the StrategyId, e.g.

```
terraform import tencentcloud_fwm_strategy.example fwms-s1d49yaf
```