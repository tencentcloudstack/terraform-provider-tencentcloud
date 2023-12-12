Provide a resource to create a VOD snapshot by time offset template.

Example Usage

```hcl
resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 130
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}
```

Import

VOD snapshot by time offset template can be imported using the id, e.g.

```
$ terraform import tencentcloud_vod_snapshot_by_time_offset_template.foo 46906
```