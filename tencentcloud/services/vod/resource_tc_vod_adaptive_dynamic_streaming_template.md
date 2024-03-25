Provide a resource to create a VOD adaptive dynamic streaming template.

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "adaptive-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = true
    tehd_config {
		  type = "TEHD-100"
	  }
  }
}
```

Import

VOD adaptive dynamic streaming template can be imported using the id($subAppId#$templateId), e.g.

```
$ terraform import tencentcloud_vod_adaptive_dynamic_streaming_template.foo $subAppId#$templateId
```