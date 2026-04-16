Provides a resource to manage an instant transcoding template for TencentCloud EdgeOne (TEO).

Example Usage

Basic Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-12345678"
  template_name = "my-transcode-template"
}
```

Advanced Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-12345678"
  template_name = "my-transcode-template"
  comment       = "My custom transcode template for edge computing"
  video_stream_switch = "on"
  audio_stream_switch = "on"

  video_template {
    codec                 = "H.264"
    fps                   = 30
    bitrate               = 2000
    resolution_adaptive   = "open"
    width                 = 1280
    height                = 720
    fill_type             = "black"
  }

  audio_template {
    codec        = "libfdk_aac"
    audio_channel = 2
  }
}
```

Video Stream Only

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "video_only" {
  zone_id       = "zone-12345678"
  template_name = "video-only-template"
  comment       = "Video transcoding only, no audio processing"
  video_stream_switch = "on"
  audio_stream_switch = "off"

  video_template {
    codec                 = "H.265"
    fps                   = 25
    bitrate               = 1500
    resolution_adaptive   = "open"
    width                 = 1920
    height                = 1080
    fill_type             = "black"
  }
}
```

Audio Stream Only

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "audio_only" {
  zone_id       = "zone-12345678"
  template_name = "audio-only-template"
  comment       = "Audio transcoding only, no video processing"
  video_stream_switch = "off"
  audio_stream_switch = "on"

  audio_template {
    codec        = "libfdk_aac"
    audio_channel = 2
  }
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID where the template is created.
* `template_name` - (Required, ForceNew) Instant transcoding template name. Length limit: 64 characters.
* `comment` - (Optional, ForceNew) Template description. Length limit: 256 characters.
* `video_stream_switch` - (Optional, ForceNew) Video stream switch. Valid values: `on` (enable), `off` (disable). Default: `on`.
* `audio_stream_switch` - (Optional, ForceNew) Audio stream switch. Valid values: `on` (enable), `off` (disable). Default: `on`.
* `video_template` - (Optional, ForceNew) Video template configuration. Required when `video_stream_switch` is `on`.
    * `codec` - (Optional) Video stream encoding format. Valid values: `H.264`, `H.265`.
    * `fps` - (Optional) Video frame rate. Range: [0, 30]. Unit: Hz. Default: 0 (same as source, max 30).
    * `bitrate` - (Optional) Video bitrate. Range: 0 or [128, 10000]. Unit: kbps. Default: 0 (automatic).
    * `resolution_adaptive` - (Optional) Resolution adaptive mode. Valid values: `open` (width/height represent long/short sides), `close` (width/height represent exact dimensions). Default: `open`.
    * `width` - (Optional) Video width (or long side) maximum. Range: 0 or [128, 1920]. Unit: px. Default: 0 (same as source).
    * `height` - (Optional) Video height (or short side) maximum. Range: 0 or [128, 1080]. Unit: px. Default: 0 (same as source).
    * `fill_type` - (Optional) Fill type for resolution mismatch. Valid values: `stretch`, `black`, `white`, `gauss`. Default: `black`.
* `audio_template` - (Optional, ForceNew) Audio template configuration. Required when `audio_stream_switch` is `on`.
    * `codec` - (Optional) Audio stream encoding format. Valid value: `libfdk_aac`.
    * `audio_channel` - (Optional) Audio channel count. Valid value: 2 (stereo). Default: 2.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID in the format `zone_id#template_id`.
* `template_id` - Instant transcoding template unique identifier.
* `create_time` - Template creation time in ISO date format.
* `update_time` - Template last update time in ISO date format.

Timeouts

The `timeouts` block allows you to specify timeouts for certain operations:

* `create` - (Default: 30 minutes) Time to wait for the template to be created.
* `delete` - (Default: 30 minutes) Time to wait for the template to be deleted.

Import

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
