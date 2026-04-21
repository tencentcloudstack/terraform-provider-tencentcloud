# tencentcloud_teo_just_in_time_transcode_template

Provides a resource to manage TEO just-in-time transcode template.

## Example Usage

### Minimal Configuration

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-1234567890"
  template_name = "my-template"
}
```

### Full Configuration

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-1234567890"
  template_name = "my-template"
  comment       = "Example transcode template"

  video_stream_switch = "on"
  audio_stream_switch = "on"

  video_template {
    video_codec          = "H.264"
    fps                  = 30
    bitrate              = 2000
    resolution_adaptive  = "open"
    width                = 1280
    height               = 720
    fill_type            = "black"
  }

  audio_template {
    codec         = "libfdk_aac"
    audio_channel = 2
  }
}
```

### Video Stream Disabled

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-1234567890"
  template_name = "audio-only-template"
  comment       = "Audio only template"

  video_stream_switch = "off"
  audio_stream_switch = "on"

  audio_template {
    codec         = "libfdk_aac"
    audio_channel = 2
  }
}
```

### Audio Stream Disabled

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-1234567890"
  template_name = "video-only-template"
  comment       = "Video only template"

  video_stream_switch = "on"
  audio_stream_switch = "off"

  video_template {
    video_codec          = "H.264"
    fps                  = 30
    bitrate              = 2000
    resolution_adaptive  = "open"
    width                = 1280
    height               = 720
    fill_type            = "black"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `template_name` - (Required, ForceNew) Transcode template name. Max length: 64 characters.
* `comment` - (Optional, ForceNew) Template description. Max length: 256 characters.
* `video_stream_switch` - (Optional, ForceNew) Video stream switch. Valid values: `on`, `off`. Default: `on`. When set to `off`, `video_template` is not required.
* `audio_stream_switch` - (Optional, ForceNew) Audio stream switch. Valid values: `on`, `off`. Default: `on`. When set to `off`, `audio_template` is not required.
* `video_template` - (Optional, ForceNew) Video stream configuration parameters. Required when `video_stream_switch` is `on`.
  * `video_codec` - (Optional) Video codec. Optional values: `H.264`, `H.265`.
  * `fps` - (Optional) Video frame rate. Range: `[0, 30]`. Default: `0`.
  * `bitrate` - (Optional) Video bitrate in kbps. Range: `0` or `[128, 10000]`. Default: `0`.
  * `resolution_adaptive` - (Optional) Resolution adaptive mode. Optional values: `open`, `close`. Default: `open`.
  * `width` - (Optional) Video width/long-edge in pixels. Range: `0` or `[128, 1920]`. Default: `0`.
  * `height` - (Optional) Video height/short-edge in pixels. Range: `0` or `[128, 1080]`. Default: `0`.
  * `fill_type` - (Optional) Fill type. Optional values: `stretch`, `black`, `white`, `gauss`. Default: `black`.
* `audio_template` - (Optional, ForceNew) Audio stream configuration parameters. Required when `audio_stream_switch` is `on`.
  * `codec` - (Optional) Audio codec. Optional values: `libfdk_aac`.
  * `audio_channel` - (Optional) Audio channel count. Optional values: `2`. Default: `2`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID in the format of `zone_id#template_id`.
* `template_id` - Template unique identifier.
* `type` - Template type. Values: `preset`, `custom`.
* `create_time` - Template creation time in ISO 8601 format.
* `update_time` - Template last update time in ISO 8601 format.

## Timeouts

The `timeouts` block allows you to specify timeouts for certain operations:

* `create` - (Default 10 minutes) Timeout for creating the template.
* `read` - (Default 10 minutes) Timeout for reading the template.
* `delete` - (Default 10 minutes) Timeout for deleting the template.

## Important Notes

### ForceNew Behavior

All parameters in this resource are marked as `ForceNew`, which means any change to these parameters will force the recreation of the resource. This is because the TEO API does not provide an Update interface for just-in-time transcode templates. When a parameter is changed, the existing template will be deleted and a new one will be created with the updated configuration.

### Asynchronous Operations

Template creation and deletion are asynchronous operations. After calling the Create or Delete API, the resource will poll the Describe API until the operation completes. This may take some time depending on the service load.

### Validation Rules

- `template_name`: Maximum 64 characters.
- `comment`: Maximum 256 characters.
- `video_stream_switch`: Must be `on` or `off`.
- `audio_stream_switch`: Must be `on` or `off`.
- When `video_stream_switch` is `on`, `video_template` is required.
- When `audio_stream_switch` is `on`, `audio_template` is required.

### Template Parameters

- `fps` (video template): Range [0, 30]. When set to 0, frame rate follows the source video with a maximum of 30.
- `bitrate` (video template): Range [128, 10000] or 0. When set to 0, bitrate is automatically selected based on video quality.
- `width` and `height` (video template):
  - Both 0: Resolution follows the source.
  - Width 0, Height non-zero: Width is scaled proportionally.
  - Width non-zero, Height 0: Height is scaled proportionally.
  - Both non-zero: Resolution uses user-specified values.
- `fill_type`: Controls how the video is filled when the target aspect ratio differs from the source:
  - `stretch`: Stretch each frame to fill the entire canvas (may distort video).
  - `black`: Keep aspect ratio and fill with black borders.
  - `white`: Keep aspect ratio and fill with white borders.
  - `gauss`: Keep aspect ratio and fill with Gaussian blur.
- `resolution_adaptive`:
  - `open`: `width` is the long-edge, `height` is the short-edge.
  - `close`: `width` is the width, `height` is the height.
- `audio_channel`: Currently only supports `2` (stereo).

## Import

TEO just-in-time transcode template can be imported using the resource ID in the format `zone_id#template_id`. For example:

```bash
terraform import tencentcloud_teo_just_in_time_transcode_template.example zone-1234567890#tpl-abcdefghij
```
