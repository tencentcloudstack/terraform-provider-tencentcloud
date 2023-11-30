Use this data source to query detailed information of eb event_rules
Example Usage
```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus_rule"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_eb_event_rule" "event_rule" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}
data "tencentcloud_eb_event_rules" "event_rules" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  order_by     = "AddTime"
  order        = "DESC"
  depends_on = [tencentcloud_eb_event_rule.event_rule]
}
```