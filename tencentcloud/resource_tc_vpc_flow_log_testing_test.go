package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingVpcFlowLogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingVpcFlowLog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_flow_log.flow_log", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_name", "iac-test-1"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_description", "this is a testing flow log"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_flow_log.flow_log",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cloud_log_region",
					"flow_log_storage",
				},
			},
			{
				Config: testAccTestingVpcFlowLogUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_flow_log.flow_log", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_name", "iac-test-2"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_description", "updated"),
				),
			},
		},
	})
}

const testAccTestingVpcFlowLog = `

resource "tencentcloud_vpc_flow_log" "flow_log" {
  flow_log_name = "iac-test-1"
  resource_type = "NETWORKINTERFACE"
  resource_id = "eni-qz9wxgmd"
  traffic_type = "ACCEPT"
  vpc_id = "vpc-humgpppd"
  flow_log_description = "this is a testing flow log"
  cloud_log_id = "e6acd27c-365c-4959-8257-751d86657439" # FIXME use data.logsets (not supported) instead
  storage_type = "cls"
}
`

const testAccTestingVpcFlowLogUpdate = `
resource "tencentcloud_vpc_flow_log" "flow_log" {
  flow_log_name = "iac-test-2"
  resource_type = "NETWORKINTERFACE"
  resource_id = "eni-qz9wxgmd"
  traffic_type = "ACCEPT"
  vpc_id = "vpc-humgpppd"
  flow_log_description = "updated"
  cloud_log_id = "e6acd27c-365c-4959-8257-751d86657439" # FIXME use data.logsets (not supported) instead
  storage_type = "cls"
}
`
