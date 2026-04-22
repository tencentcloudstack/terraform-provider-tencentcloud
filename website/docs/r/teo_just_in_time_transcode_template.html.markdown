---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_just_in_time_transcode_template"
sidebar_current: "docs-tencentcloud-resource-teo_just_in_time_transcode_template"
description: |-
  Provides a resource to create a TEO just-in-time transcode template.
---

# tencentcloud_teo_just_in_time_transcode_template

Provides a resource to create a TEO just-in-time transcode template.

## Example Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-2qtuhspy7cr6"
  template_name = "my-template"
  comment       = "Example transcode template"

  video_stream_switch = "on"
  audio_stream_switch = "on"

  video_template {
    video_codec         = "H.264"
    fps                 = 30
    bitrate             = 2000
    resolution_adaptive = "open"
    width               = 1280
    height              = 720
    fill_type           = "black"
  }

  audio_template {
    codec         = "libfdk_aac"
    audio_channel = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String, ForceNew) Transcode template name. Max length: 64 characters.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `audio_stream_switch` - (Optional, String) Audio stream switch. Valid values: on, off. Default: on.
* `audio_template` - (Optional, List) Audio stream configuration parameters. Required when audio_stream_switch is on.
* `comment` - (Optional, String) Template description. Max length: 256 characters.
* `video_stream_switch` - (Optional, String) Video stream switch. Valid values: on, off. Default: on.
* `video_template` - (Optional, List) Video stream configuration parameters. Required when video_stream_switch is on.

The `audio_template` object supports the following:

* `audio_channel` - (Optional, Int) Audio channel count. Optional values: 2. Default: 2.
* `codec` - (Optional, String) Audio codec. Optional values: libfdk_aac.

The `video_template` object supports the following:

* `bitrate` - (Optional, Int) Video bitrate in kbps. Range: 0 or [128, 10000]. Default: 0.
* `fill_type` - (Optional, String) Fill type. Optional values: stretch, black, white, gauss. Default: black.
* `fps` - (Optional, Float64) Video frame rate. Range: [0, 30]. Default: 0.
* `height` - (Optional, Int) Video height/short-edge in pixels. Range: 0 or [128, 1080]. Default: 0.
* `resolution_adaptive` - (Optional, String) Resolution adaptive mode. Optional values: open, close. Default: open.
* `video_codec` - (Optional, String) Video codec. Optional values: H.264, H.265.
* `width` - (Optional, Int) Video width/long-edge in pixels. Range: 0 or [128, 1920]. Default: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Template creation time in ISO 8601 format.
* `template_id` - Template ID returned after creation.
* `type` - Template type. Values: preset, custom.
* `update_time` - Template last update time in ISO 8601 format.


## Import

TEO just-in-time transcode template can be imported using the joint id "zone_id#template_id", e.g.

```
terraform import tencentcloud_teo_just_in_time_transcode_template.example zone-2qtuhspy7cr6#jitt-abcdefghij
```

