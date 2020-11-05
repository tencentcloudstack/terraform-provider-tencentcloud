---
subcategory: "VOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_super_player_configs"
sidebar_current: "docs-tencentcloud-datasource-vod_super_player_configs"
description: |-
  Use this data source to query detailed information of VOD super player configs.
---

# tencentcloud_vod_super_player_configs

Use this data source to query detailed information of VOD super player configs.

## Example Usage

```hcl
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

data "tencentcloud_vod_super_player_configs" "foo" {
  type = "Custom"
  name = "tf-super-player"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of super player config.
* `result_output_file` - (Optional) Used to save results.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.
* `type` - (Optional) Config type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config_list` - A list of super player configs. Each element contains the following attributes:
  * `adaptive_dynamic_streaming_definition` - ID of the unencrypted adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `false`.
  * `comment` - Template description.
  * `create_time` - Creation time of template in ISO date format.
  * `domain` - Domain name used for playback. If it is left empty or set to `Default`, the domain name configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used.
  * `drm_streaming_info` - Content of the DRM-protected adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `true`.
    * `simple_aes_definition` - ID of the adaptive dynamic streaming template whose protection type is `SimpleAES`.
  * `drm_switch` - Switch of DRM-protected adaptive bitstream playback: `true`: enabled, indicating to play back only output adaptive bitstreams protected by DRM; `false`: disabled, indicating to play back unencrypted output adaptive bitstreams.
  * `image_sprite_definition` - ID of the image sprite template that allows output.
  * `name` - Player configuration name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.
  * `resolution_names` - Display name of player for substreams with different resolutions. If this parameter is left empty or an empty array, the default configuration will be used: `min_edge_length: 240, name: LD`; `min_edge_length: 480, name: SD`; `min_edge_length: 720, name: HD`; `min_edge_length: 1080, name: FHD`; `min_edge_length: 1440, name: 2K`; `min_edge_length: 2160, name: 4K`; `min_edge_length: 4320, name: 8K`.
    * `min_edge_length` - Length of video short side in px.
    * `name` - Display name.
  * `scheme` - Scheme used for playback. If it is left empty or set to `Default`, the scheme configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. Other valid values: `HTTP`; `HTTPS`.
  * `type` - Template type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.
  * `update_time` - Last modified time of template in ISO date format.


