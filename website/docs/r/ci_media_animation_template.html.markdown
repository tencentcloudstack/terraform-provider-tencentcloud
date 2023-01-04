---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_animation_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_animation_template"
description: |-
  Provides a resource to create a ci media_animation_template
---

# tencentcloud_ci_media_animation_template

Provides a resource to create a ci media_animation_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_animation_template" "media_animation_template" {
  bucket = "terraform-ci-1308919341"
  name   = "animation_template-002"
  container {
    format = "gif"
  }
  video {
    codec                          = "gif"
    width                          = "1280"
    height                         = ""
    fps                            = "20"
    animate_only_keep_key_frame    = "true"
    animate_time_interval_of_frame = ""
    animate_frames_per_second      = ""
    quality                        = ""

  }
  time_interval {
    start    = "0"
    duration = "60"

  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `container` - (Required, List) container format.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `time_interval` - (Optional, List) time interval.
* `video` - (Optional, List) video information, do not upload Video, which is equivalent to deleting video information.

The `container` object supports the following:

* `format` - (Required, String) Package format.

The `time_interval` object supports the following:

* `duration` - (Optional, String) duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.
* `start` - (Optional, String) Starting time, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.

The `video` object supports the following:

* `codec` - (Required, String) Codec format `gif`, `webp`.
* `animate_frames_per_second` - (Optional, String) Animation per second frame number, Priority: AnimateFramesPerSecond &gt; AnimateOnlyKeepKeyFrame &gt; AnimateTimeIntervalOfFrame.
* `animate_only_keep_key_frame` - (Optional, String) GIFs are kept only Keyframe, Priority: AnimateFramesPerSecond &gt; AnimateOnlyKeepKeyFrame &gt; AnimateTimeIntervalOfFrame.
* `animate_time_interval_of_frame` - (Optional, String) Animation frame extraction every time, (0, video duration], Animation frame extraction time interval, If TimeInterval.Duration is set, it is less than this value.
* `fps` - (Optional, String) Frame rate, value range: (0, 60], Unit: fps.
* `height` - (Optional, String) High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video, must be even.
* `quality` - (Optional, String) Set relative quality, [1, 100), webp image quality setting takes effect, gif has no quality parameter.
* `width` - (Optional, String) width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_animation_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_animation_template.media_animation_template terraform-ci-xxxxxx#t18210645f96564eaf80e86b1f58c20152
```

