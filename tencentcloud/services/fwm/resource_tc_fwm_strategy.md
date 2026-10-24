Provides a resource to create a FWM (Firewall Manager) strategy.

Example Usage

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product        = "enterprise_sg"
  sub_account_uin = "100028274628"
  exec_area      = "pre"
  group_id      = "fwmg-xxxxxx"
  sequence       = 1
  receive_account = ["100028274628"]
}
```

Import

FWM strategy can be imported using the StrategyId, e.g.

```
terraform import tencentcloud_fwm_strategy.example fwms-xxxx
```