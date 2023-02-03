---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_snapshot_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_snapshot_template"
description: |-
  Provides a resource to create a ci media_snapshot_template
---

# tencentcloud_ci_media_snapshot_template

Provides a resource to create a ci media_snapshot_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_snapshot_template" "media_snapshot_template" {
  bucket = "terraform-ci-xxxxxx"
  name   = "snapshot_template_test"
  snapshot {
    count             = "10"
    snapshot_out_mode = "SnapshotAndSprite"
    sprite_snapshot_config {
      color   = "White"
      columns = "10"
      lines   = "10"
      margin  = "10"
      padding = "10"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `snapshot` - (Required, List) screenshot.

The `snapshot` object supports the following:

* `count` - (Required, String) Number of screenshots, range (0 10000].
* `black_level` - (Optional, String) Screenshot black screen detection parameters, Valid when IsCheckBlack=true, Value reference range [30, 100], indicating the proportion of black pixels, the smaller the value, the smaller the proportion of black pixels, Start&gt;0, the parameter setting is invalid, no filter black screen, Start =0 parameter is valid, the start time of the frame capture is the first frame non-black screen start.
* `ci_param` - (Optional, String) Screenshot image processing parameters, for example: imageMogr2/format/png.
* `height` - (Optional, String) high, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video.
* `is_check_black` - (Optional, String) Whether to enable black screen detection true/false.
* `is_check_count` - (Optional, String) Whether to check the number of screenshots forcibly, when using custom interval mode to take screenshots, the video time is not long enough to capture Count screenshots, you can switch to average screenshot mode to capture Count screenshots.
* `mode` - (Optional, String) Screenshot mode, value range: {Interval, Average, KeyFrame}- Interval means interval mode Average means average mode- KeyFrame represents the key frame mode- Interval mode: Start, TimeInterval, The Count parameter takes effect. When Count is set and TimeInterval is not set, Indicates to capture all frames, a total of Count pictures- Average mode: Start, the Count parameter takes effect. express.
* `pixel_black_threshold` - (Optional, String) Screenshot black screen detection parameters, Valid when IsCheckBlack=true, The threshold for judging whether a pixel is a black point, value range: [0, 255].
* `snapshot_out_mode` - (Optional, String) Screenshot output mode parameters, Value range: {OnlySnapshot, OnlySprite, SnapshotAndSprite}, OnlySnapshot means output only screenshot mode OnlySprite means only output sprite mode SnapshotAndSprite means output screenshot and sprite mode.
* `sprite_snapshot_config` - (Optional, List) Screenshot output configuration.
* `start` - (Optional, String) Starting time, [0 video duration] in seconds, Support float format, the execution accuracy is accurate to milliseconds.
* `time_interval` - (Optional, String) Screenshot time interval, (0 3600], in seconds, Support float format, the execution accuracy is accurate to milliseconds.
* `width` - (Optional, String) wide, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video.

The `sprite_snapshot_config` object supports the following:

* `color` - (Required, String) See `https://www.ffmpeg.org/ffmpeg-utils.html#color-syntax` for details on supported colors.
* `columns` - (Required, String) Number of screenshot columns, value range: [1, 10000].
* `lines` - (Required, String) Number of screenshot lines, value range: [1, 10000].
* `cell_height` - (Optional, String) Single image height Value range: [8, 4096], Unit: px.
* `cell_width` - (Optional, String) Single image width Value range: [8, 4096], Unit: px.
* `margin` - (Optional, String) screenshot margin size, Value range: [8, 4096], Unit: px.
* `padding` - (Optional, String) screenshot padding size, Value range: [8, 4096], Unit: px.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - creation time.
* `template_id` - Template ID.
* `update_time` - update time.


## Import

ci media_snapshot_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_snapshot_template.media_snapshot_template terraform-ci-xxxxxx#t18210645f96564eaf80e86b1f58c20152
```

