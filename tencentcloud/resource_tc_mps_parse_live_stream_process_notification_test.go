package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsParseLiveStreamProcessNotificationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsParseLiveStreamProcessNotification,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_parse_live_stream_process_notification.parse_live_stream_process_notification", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_parse_live_stream_process_notification.parse_live_stream_process_notification",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsParseLiveStreamProcessNotification = `

resource "tencentcloud_mps_parse_live_stream_process_notification" "parse_live_stream_process_notification" {
  content = &lt;nil&gt;
}

`
