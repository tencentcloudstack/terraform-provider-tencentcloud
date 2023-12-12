Provides a resource to create a mps sample_snapshot_template

Example Usage

```hcl
resource "tencentcloud_mps_sample_snapshot_template" "sample_snapshot_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-test-for"
  resolution_adaptive = "open"
  sample_interval     = 10
  sample_type         = "Percent"
  width               = 140
}
```

Import

mps sample_snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_sample_snapshot_template.sample_snapshot_template sample_snapshot_template_id
```