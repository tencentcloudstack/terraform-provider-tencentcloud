package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVodSuperPlayerConfigs(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSuperPlayerConfigs,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vod_super_player_configs.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.name", "tf-super-player1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.drm_switch", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.drm_streaming_info.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.resolution_names.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.resolution_names.0.min_edge_length", "889"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.resolution_names.0.name", "test1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.resolution_names.1.min_edge_length", "890"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.resolution_names.1.name", "test2"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.domain", "Default"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.scheme", "Default"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.comment", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.drm_streaming_info.0.simple_aes_definition"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.image_sprite_definition"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.update_time"),
				),
			},
		},
	})
}

const testAccDataSourceVodSuperPlayerConfig = testAccVodAdaptiveDynamicStreamingTemplate + testAccVodImageSpriteTemplate + `
resource "tencentcloud_vod_super_player_config" "foo" {
  name                    = "tf-super-player1"
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
`

const testAccVodSuperPlayerConfigs = testAccDataSourceVodSuperPlayerConfig + `
data "tencentcloud_vod_super_player_configs" "foo" {
  type = "Custom"
  name = tencentcloud_vod_super_player_config.foo.id
}
`
