package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPtsAlertChannel_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsAlertChannel,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_alert_channel.alert_channel", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_alert_channel.alertChannel",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsAlertChannel = `

resource "tencentcloud_pts_alert_channel" "alert_channel" {
  notice_id = ""
  project_id = ""
  a_m_p_consumer_id = ""
            }

`
