Provides a resource to create a eb event_bus

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
```

Import

eb event_bus can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_bus.event_bus event_bus_id
```