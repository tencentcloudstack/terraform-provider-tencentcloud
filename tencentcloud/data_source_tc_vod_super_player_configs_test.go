package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudVodSuperPlayerConfigs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSuperPlayerConfigs,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vod_super_player_configs.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.name", "tf-super-player"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_super_player_configs.foo", "config_list.0.drm_switch", "ON"),
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

const testAccVodSuperPlayerConfigs = testAccVodSuperPlayerConfig + `
data "tencentcloud_vod_super_player_configs" "foo" {
  type = "Custom"
  name = tencentcloud_vod_super_player_config.foo.id
}
`
