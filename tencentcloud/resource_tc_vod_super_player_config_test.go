package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_vod_super_player_config
	resource.AddTestSweepers("tencentcloud_vod_super_player_config", &resource.Sweeper{
		Name: "tencentcloud_vod_super_player_config",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			sharedClient, err := sharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(*TencentCloudClient)
			vodService := VodService{
				client: client.apiV3Conn,
			}
			filter := make(map[string]interface{})
			configs, e := vodService.DescribeSuperPlayerConfigsByFilter(ctx, filter)
			if e != nil {
				return nil
			}
			for _, config := range configs {
				ee := vodService.DeleteSuperPlayerConfig(ctx, *config.Name, uint64(0))
				if ee != nil {
					continue
				}
			}
			return nil
		},
	})
}
func TestAccTencentCloudVodSuperPlayerConfigResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVodSuperPlayerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSuperPlayerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodSuperPlayerConfigExists("tencentcloud_vod_super_player_config.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "name", "tf-super-player-0"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "drm_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "drm_streaming_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.0.min_edge_length", "889"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.0.name", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.1.min_edge_length", "890"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.1.name", "test2"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "domain", "Default"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "scheme", "Default"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "comment", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_super_player_config.foo", "drm_streaming_info.0.simple_aes_definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_super_player_config.foo", "image_sprite_definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_super_player_config.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_super_player_config.foo", "update_time"),
				),
			},
			{
				Config: testAccVodSuperPlayerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "drm_switch", "false"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.0.min_edge_length", "891"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.0.name", "test1-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.1.min_edge_length", "892"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "resolution_names.1.name", "test2-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_super_player_config.foo", "comment", "test-update"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_super_player_config.foo", "adaptive_dynamic_streaming_definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_super_player_config.foo", "image_sprite_definition"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_super_player_config.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sub_app_id"},
			},
		},
	})
}

func testAccCheckVodSuperPlayerConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	vodService := VodService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_super_player_config" {
			continue
		}

		_, has, err := vodService.DescribeSuperPlayerConfigsById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("vod super player config still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVodSuperPlayerConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod super player config %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod super player config id is not set")
		}
		vodService := VodService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, has, err := vodService.DescribeSuperPlayerConfigsById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("vod super player config doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodSuperPlayerConfig = testAccVodAdaptiveDynamicStreamingTemplate + testAccVodImageSpriteTemplate + `
resource "tencentcloud_vod_super_player_config" "foo" {
  name                    = "tf-super-player-0"
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

const testAccVodSuperPlayerConfigUpdate = testAccVodAdaptiveDynamicStreamingTemplate + testAccVodImageSpriteTemplate + `
resource "tencentcloud_vod_super_player_config" "foo" {
  name                                  = "tf-super-player-0"
  drm_switch                            = false
  adaptive_dynamic_streaming_definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
  image_sprite_definition               = tencentcloud_vod_image_sprite_template.foo.id
  resolution_names {
    min_edge_length = 891
    name            = "test1-update"
  }
  resolution_names {
    min_edge_length = 892
    name            = "test2-update"
  }
  domain                                = "Default"
  scheme                                = "Default"
  comment                               = "test-update"
}
`
