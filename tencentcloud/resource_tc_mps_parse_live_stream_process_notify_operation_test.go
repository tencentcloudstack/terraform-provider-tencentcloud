package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsParseLiveStreamProcessNotifyOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsParseLiveStreamProcessNotifyOperation,
				Check:  resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_parse_live_stream_process_notify_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_parse_live_stream_process_notify_operation.operation", "content"),
				),
			},
		},
	})
}

const testAccMpsParseLiveStreamProcessNotifyOperation = `

resource "tencentcloud_mps_parse_live_stream_process_notify_operation" "operation" {
  content = ""
}

`
