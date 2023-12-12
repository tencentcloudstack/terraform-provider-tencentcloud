Provides a resource to create a ci media_snapshot_template

Example Usage

```hcl
resource "tencentcloud_ci_media_snapshot_template" "media_snapshot_template" {
    bucket = "terraform-ci-xxxxxx"
  	name = "snapshot_template_test"
  	snapshot {
      count = "10"
      snapshot_out_mode = "SnapshotAndSprite"
      sprite_snapshot_config {
        color = "White"
        columns = "10"
        lines = "10"
        margin = "10"
        padding = "10"
      }
  	}
}
```

Import

ci media_snapshot_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_snapshot_template.media_snapshot_template terraform-ci-xxxxxx#t18210645f96564eaf80e86b1f58c20152
```