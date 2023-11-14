package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsParseLiveStreamProcessNotificationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsParseLiveStreamProcessNotificationDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_parse_live_stream_process_notification.parse_live_stream_process_notification")),
			},
		},
	})
}

const testAccMpsParseLiveStreamProcessNotificationDataSource = `

data "tencentcloud_mps_parse_live_stream_process_notification" "parse_live_stream_process_notification" {
  content = ""
}

`
