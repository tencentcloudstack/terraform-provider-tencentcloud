Provides a resource to create events audit track

Example Usage

```hcl
resource "tencentcloud_events_audit_track" "example" {
  name = "track_example"

  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "393953ac-5c1b-457d-911d-376271b1b4f2"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cls"
  }

  filters {
    resource_fields {
      resource_type = "cam"
      action_type   = "*"
      event_names   = ["AddSubAccount", "AddSubAccountCheckingMFA"]
    }
    resource_fields {
      resource_type = "cvm"
      action_type   = "*"
      event_names   = ["*"]
    }
    resource_fields {
      resource_type = "tke"
      action_type   = "*"
      event_names   = ["*"]
    }
  }
}
```

Import

events audit track can be imported using the id, e.g.
```
$ terraform import tencentcloud_events_audit_track.example 24283
```