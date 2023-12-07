Provides a resource to create a mps event

Example Usage

```hcl
resource "tencentcloud_mps_event" "event" {
  event_name = "you-event-name"
  description = "event description"
}
```

Import

mps event can be imported using the id, e.g.

```
terraform import tencentcloud_mps_event.event event_id
```