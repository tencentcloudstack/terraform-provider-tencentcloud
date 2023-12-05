Use this data source to query detailed information of VOD super player configs.

Example Usage

```hcl
resource "tencentcloud_vod_super_player_config" "foo" {
  name                    = "tf-super-player"
  drm_switch              = true
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
  domain                  = "Default"
  scheme                  = "Default"
  comment                 = "test"
}

data "tencentcloud_vod_super_player_configs" "foo" {
  type = "Custom"
  name = "tf-super-player"
}
```