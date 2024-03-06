Provides a resource to create a vod snapshot template

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "snapshotTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_sample_snapshot_template" "sample_snapshot_template" {
  sample_type = "Percent"
  sample_interval = 10
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name = "testSampleSnapshot"
  width = 500
  height = 400
  resolution_adaptive = "open"
  format = "jpg"
  comment = "test sample snopshot"
  fill_type = "black"
}
```

Import

vod snapshot template can be imported using the id, e.g.

```
terraform import tencentcloud_vod_sample_snapshot_template.sample_snapshot_template $subAppId#$templateId
```