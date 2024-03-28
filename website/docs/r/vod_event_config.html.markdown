---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_event_config"
sidebar_current: "docs-tencentcloud-resource-vod_event_config"
description: |-
  Provide a resource to create a vod event config.
---

# tencentcloud_vod_event_config

Provide a resource to create a vod event config.

## Example Usage

```hcl
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "eventconfig-subapplication"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_event_config" "event_config" {
  mode                               = "PUSH"
  notification_url                   = "http://mydemo.com/receiveevent"
  upload_media_complete_event_switch = "ON"
  delete_media_complete_event_switch = "ON"
  sub_app_id                         = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
}
```

## Argument Reference

The following arguments are supported:

* `sub_app_id` - (Required, Int) Sub app id.
* `delete_media_complete_event_switch` - (Optional, String) Whether to receive video deletion completion event notification, default `OFF` is to ignore the event notification, `ON` is to receive event notification.
* `mode` - (Optional, String) How to receive event notifications. Valid values:
- Push: HTTP callback notification;
- PULL: Reliable notification based on message queuing.
* `notification_url` - (Optional, String) The address used to receive 3.0 format callbacks when receiving HTTP callback notifications. Note: If you take the NotificationUrl parameter and the value is an empty string, the 3.0 format callback address is cleared.
* `upload_media_complete_event_switch` - (Optional, String) Whether to receive video upload completion event notification, default `OFF` means to ignore the event notification, `ON` means to receive event notification.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VOD event config can be imported using the subAppId, e.g.

```
$ terraform import tencentcloud_vod_event_config.foo $subAppId
```

