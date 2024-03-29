package vod_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVodEventConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodEventConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vod_event_config.event_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_event_config.event_config", "sub_app_id"),
					resource.TestCheckResourceAttr("tencentcloud_vod_event_config.event_config", "mode", "PUSH"),
					resource.TestCheckResourceAttr("tencentcloud_vod_event_config.event_config", "notification_url", "http://mydemo.com/receiveevent"),
					resource.TestCheckResourceAttr("tencentcloud_vod_event_config.event_config", "upload_media_complete_event_switch", "OFF"),
					resource.TestCheckResourceAttr("tencentcloud_vod_event_config.event_config", "delete_media_complete_event_switch", "OFF"),
				),
			},
			{
				Config: testAccVodEventConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_event_config.event_config", "upload_media_complete_event_switch", "ON"),
					resource.TestCheckResourceAttr("tencentcloud_vod_event_config.event_config", "delete_media_complete_event_switch", "ON"),
				),
			},
			{
				ResourceName:      "tencentcloud_vod_event_config.event_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVodEventConfig = `

resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "eventconfig-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_event_config" "event_config" {
  mode = "PUSH"
  notification_url = "http://mydemo.com/receiveevent"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
}

`

const testAccVodEventConfigUpdate = `

resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "eventconfig-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_event_config" "event_config" {
  mode = "PUSH"
  notification_url = "http://mydemo.com/receiveevent"
  upload_media_complete_event_switch = "ON"
  delete_media_complete_event_switch = "ON"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
}

`
