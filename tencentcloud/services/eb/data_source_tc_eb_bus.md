Use this data source to query detailed information of eb bus

Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}
data "tencentcloud_eb_bus" "bus" {
  order_by = "AddTime"
  order = "DESC"
  filters {
		values = ["Custom"]
		name = "Type"
  }

  depends_on = [ tencentcloud_eb_event_bus.foo ]
}
```