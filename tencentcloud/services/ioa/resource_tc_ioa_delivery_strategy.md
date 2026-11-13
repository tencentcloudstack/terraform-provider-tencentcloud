Provides a resource to create or edit an IOA (T-Sec Internet of Amount Security Management Center) log delivery strategy.

Example Usage

```hcl
resource "tencentcloud_ioa_delivery_strategy" "example" {
  name        = "tf-example"
  description = "delivery strategy for IOA logs"
  type        = 2
  config_id   = 944
  enabled     = true
  log_types = [
    "040201000000",
    "040205000000",
    "040209000000",
    "320500000000",
    "340600000000",
    "340100000000",
    "040206000000",
    "040203000000",
    "040202000000",
    "040207000000",
    "320300000000",
    "340200000000"
  ]
}
```

Import

IOA delivery strategy can be imported using the id, e.g.

```
terraform import tencentcloud_ioa_delivery_strategy.example 621
```
