---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_just_in_time_transcode_template"
sidebar_current: "docs-tencentcloud-resource-teo_just_in_time_transcode_template"
description: |-
  Provides a resource to manage an instant transcoding template for TencentCloud EdgeOne (TEO).
---

# tencentcloud_teo_just_in_time_transcode_template

Provides a resource to manage an instant transcoding template for TencentCloud EdgeOne (TEO).

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-12345678"
  template_name = "my-transcode-template"
}
```

### Advanced Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id             = "zone-12345678"
  template_name       = "my-transcode-template"
  comment             = "My custom transcode template for edge computing"
  video_stream_switch = "on"
  audio_stream_switch = "on"

  video_template {
    codec               = "H.264"
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

### Video Stream Only

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "video_only" {
  zone_id             = "zone-12345678"
  template_name       = "video-only-template"
  comment             = "Video transcoding only, no audio processing"
  video_stream_switch = "on"
  audio_stream_switch = "off"

  video_template {
    codec               = "H.265"
    fps                 = 25
    bitrate             = 1500
    resolution_adaptive = "open"
    width               = 1920
    height              = 1080
    fill_type           = "black"
  }
}
```

### Audio Stream Only

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "audio_only" {
  zone_id             = "zone-12345678"
  template_name       = "audio-only-template"
  comment             = "Audio transcoding only, no video processing"
  video_stream_switch = "off"
  audio_stream_switch = "on"

  audio_template {
    codec         = "libfdk_aac"
    audio_channel = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String, ForceNew) Instant transcoding template name. Length limit: 64 characters.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `audio_stream_switch` - (Optional, String, ForceNew) Audio stream switch. Values: `on` (enable), `off` (disable). Default: `on`.
* `audio_template` - (Optional, List, ForceNew) Audio template configuration. Required when `audio_stream_switch` is `on`.
* `comment` - (Optional, String, ForceNew) Template description. Length limit: 256 characters.
* `video_stream_switch` - (Optional, String, ForceNew) Video stream switch. Values: `on` (enable), `off` (disable). Default: `on`.
* `video_template` - (Optional, List, ForceNew) Video template configuration. Required when `video_stream_switch` is `on`.

The `audio_template` object supports the following:

* `audio_channel` - (Optional, Int) Audio channel count. Value: 2 (stereo). Default: 2.
* `codec` - (Optional, String) Audio stream encoding format. Values: `libfdk_aac`.

The `video_template` object supports the following:

* `bitrate` - (Optional, Int) Video bitrate. Range: 0 or [128, 10000]. Unit: kbps. Default: 0.
* `codec` - (Optional, String) Video stream encoding format. Values: `H.264`, `H.265`.
* `fill_type` - (Optional, String) Fill type for resolution mismatch. Values: `stretch`, `black`, `white`, `gauss`. Default: `black`.
* `fps` - (Optional, Float64) Video frame rate. Range: [0, 30]. Unit: Hz. Default: 0.
* `height` - (Optional, Int) Video height (or short side) maximum. Range: 0 or [128, 1080]. Unit: px. Default: 0.
* `resolution_adaptive` - (Optional, String) Resolution adaptive mode. Values: `open` (width/height represent long/short sides), `close` (width/height represent exact dimensions). Default: `open`.
* `width` - (Optional, Int) Video width (or long side) maximum. Range: 0 or [128, 1920]. Unit: px. Default: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Template creation time in ISO date format.
* `template_id` - Instant transcoding template unique identifier.
* `update_time` - Template last update time in ISO date format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.

## Import

Just-in-time transcode template can be imported using the resource ID, which is in the format `zone_id#template_id`:

```hcl
terraform import tencentcloud_teo_just_in_time_transcode_template.example zone-abc123#tpl-def456
```

Notes

* All parameters are ForceNew, meaning any parameter change will trigger resource recreation (delete + create).
* When `video_stream_switch` is set to `on`, the `video_template` block must be provided.
* When `audio_stream_switch` is set to `on`, the `audio_template` block must be provided.
* When `video_stream_switch` is set to `off`, the `video_template` block should not be provided.
* When `audio_stream_switch` is set to `off`, the `audio_template` block should not be provided.
* The template creation and deletion operations are asynchronous. The provider will poll until the operation completes.
* In case of slow operations, consider increasing the timeout values in the `timeouts` block.

