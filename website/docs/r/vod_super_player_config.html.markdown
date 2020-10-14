---
subcategory: "VOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_super_player_config"
sidebar_current: "docs-tencentcloud-resource-vod_super_player_config"
description: |-
  Provide a resource to create a VOD super player config.
---

# tencentcloud_vod_super_player_config

Provide a resource to create a VOD super player config.

## Example Usage

```hcl
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec               = "libx265"
      fps                 = 4
      bitrate             = 129
      resolution_adaptive = false
      width               = 128
      height              = 128
      fill_type           = "stretch"
    }
    audio {
      codec         = "libmp3lame"
      bitrate       = 129
      sample_rate   = 44100
      audio_channel = "dual"
    }
    remove_audio = false
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
  }
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}

resource "tencentcloud_vod_super_player_config" "foo" {
  name       = "tf-super-player"
  drm_switch = true
  drm_streaming_info {
    simple_aes_definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
  }
  image_sprite_definition = tencentcloud_vod_image_sprite_template.foo.id
  resolution_names {
    min_edge_length = 889
    name            = "test1"
  }
  resolution_names {
    min_edge_length = 890
    name            = "test2"
  }
  domain  = "Default"
  scheme  = "Default"
  comment = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Player configuration name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.
* `adaptive_dynamic_streaming_definition` - (Optional) ID of the unencrypted adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `false`.
* `comment` - (Optional) Template description. Length limit: 256 characters.
* `domain` - (Optional) Domain name used for playback. If it is left empty or set to `Default`, the domain name configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. `Default` by default.
* `drm_streaming_info` - (Optional) Content of the DRM-protected adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `true`.
* `drm_switch` - (Optional) Switch of DRM-protected adaptive bitstream playback: `true`: enabled, indicating to play back only output adaptive bitstreams protected by DRM; `false`: disabled, indicating to play back unencrypted output adaptive bitstreams. Default value: `false`.
* `image_sprite_definition` - (Optional) ID of the image sprite template that allows output.
* `resolution_names` - (Optional) Display name of player for substreams with different resolutions. If this parameter is left empty or an empty array, the default configuration will be used: `min_edge_length: 240, name: LD`; `min_edge_length: 480, name: SD`; `min_edge_length: 720, name: HD`; `min_edge_length: 1080, name: FHD`; `min_edge_length: 1440, name: 2K`; `min_edge_length: 2160, name: 4K`; `min_edge_length: 4320, name: 8K`.
* `scheme` - (Optional) Scheme used for playback. If it is left empty or set to `Default`, the scheme configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. Other valid values: `HTTP`; `HTTPS`.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.

The `drm_streaming_info` object supports the following:

* `simple_aes_definition` - (Optional) ID of the adaptive dynamic streaming template whose protection type is `SimpleAES`.

The `resolution_names` object supports the following:

* `min_edge_length` - (Required) Length of video short side in px.
* `name` - (Required) Display name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of template in ISO date format.
* `update_time` - Last modified time of template in ISO date format.


## Import

VOD super player config can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_super_player_config.foo tf-super-player
```

