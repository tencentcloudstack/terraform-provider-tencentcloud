package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMpsStartFlowOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsStartFlowOperation_start,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "flow_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_start_flow_operation.operation", "start", "true"),
				),
			},
			{
				Config: testAccMpsStartFlowOperation_stop,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "flow_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_start_flow_operation.operation", "start", "false"),
				),
			},
		},
	})
}

const testAccMpsStartFlowOperation_basic = `
resource "tencentcloud_mps_event" "event_rtp" {
	event_name = "tf_test_event_rtp_oper"
	description = "tf test mps event description"
  }

resource "tencentcloud_mps_flow" "flow_rtp" {
  flow_name = "tf_test_flow_rtp_oper"
  max_bandwidth = 10000000
  input_group {
		input_name = "test_inputname"
		protocol = "RTP"
		description = "input name Description"
		allow_ip_list = ["0.0.0.0/0"]
		rtp_settings {
			fec = "none"
			idle_timeout = 1000
		}
  }
  event_id = tencentcloud_mps_event.event_rtp.id
}
`

const testAccMpsStartFlowOperation_start = testAccMpsStartFlowOperation_basic + `
resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id = tencentcloud_mps_flow.flow_rtp.id
  start   = true
}

`

const testAccMpsStartFlowOperation_stop = testAccMpsStartFlowOperation_basic + `
resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id = tencentcloud_mps_flow.flow_rtp.id
  start   = false
}

`
