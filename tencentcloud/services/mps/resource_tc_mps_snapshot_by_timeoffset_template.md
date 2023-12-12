Provides a resource to create a mps snapshot_by_timeoffset_template

Example Usage

```hcl
resource "tencentcloud_mps_snapshot_by_timeoffset_template" "snapshot_by_timeoffset_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-test"
  resolution_adaptive = "open"
  width               = 140
}
```

Import

mps snapshot_by_timeoffset_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_snapshot_by_timeoffset_template.snapshot_by_timeoffset_template snapshot_by_timeoffset_template_id
```