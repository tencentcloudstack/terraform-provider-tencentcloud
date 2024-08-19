Provides a resource to create a CDC site

Example Usage

Create a basic CDC site

```hcl
resource "tencentcloud_cdc_site" "example" {
  name         = "tf-example"
  country      = "China"
  province     = "Guangdong Province"
  city         = "Guangzhou"
  address_line = "Tencent Building"
  description  = "desc."
}
```

Create a complete CDC site

```hcl
resource "tencentcloud_cdc_site" "example" {
  name                  = "tf-example"
  country               = "China"
  province              = "Guangdong Province"
  city                  = "Guangzhou"
  address_line          = "Shenzhen Tencent Building"
  optional_address_line = "Shenzhen Tencent Building of Binhai"
  description           = "desc."
  fiber_type            = "MM"
  optical_standard      = "MM"
  power_connectors      = "380VAC3P"
  power_feed_drop       = "DOWN"
  max_weight            = 100
  power_draw_kva        = 10
  uplink_speed_gbps     = 10
  uplink_count          = 2
  condition_requirement = true
  dimension_requirement = true
  redundant_networking  = true
  need_help             = true
  redundant_power       = true
  breaker_requirement   = true
}
```

Import

CDC site can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_site.example site-43qcf1ag
```
