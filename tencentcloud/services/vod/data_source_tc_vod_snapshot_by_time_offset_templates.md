Use this data source to query detailed information of VOD snapshot by time offset templates.

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

data "tencentcloud_vod_snapshot_by_time_offset_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
}
```