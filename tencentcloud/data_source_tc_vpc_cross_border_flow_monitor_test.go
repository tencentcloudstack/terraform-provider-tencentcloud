package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcCrossBorderFlowMonitorDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCrossBorderFlowMonitorDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_cross_border_flow_monitor.cross_border_flow_monitor")),
			},
		},
	})
}

const testAccVpcCrossBorderFlowMonitorDataSource = `

data "tencentcloud_vpc_cross_border_flow_monitor" "cross_border_flow_monitor" {
  source_region = "ap-guangzhou"
  destination_region = "ap-singapore"
  ccn_id = "ccn-qd6z2ld1"
  ccn_uin = "979137"
  period = 60
  start_time = "2023-01-01 00:00:00"
  end_time = "2023-01-01 01:00:00"
}

`
