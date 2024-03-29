Provide a resource to create a vod event config.

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "eventconfig-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_event_config" "event_config" {
  mode = "PUSH"
  notification_url = "http://mydemo.com/receiveevent"
  upload_media_complete_event_switch = "ON"
  delete_media_complete_event_switch = "ON"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
}
```

Import

VOD event config can be imported using the subAppId, e.g.

```
$ terraform import tencentcloud_vod_event_config.foo $subAppId
```