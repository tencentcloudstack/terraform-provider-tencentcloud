package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsAlertChannelResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsAlertChannel,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_alert_channel.alert_channel", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_alert_channel.alert_channel",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsAlertChannel = `

resource "tencentcloud_pts_alert_channel" "alert_channel" {
  notice_id = &lt;nil&gt;
  project_id = &lt;nil&gt;
  a_m_p_consumer_id = &lt;nil&gt;
            }

`
