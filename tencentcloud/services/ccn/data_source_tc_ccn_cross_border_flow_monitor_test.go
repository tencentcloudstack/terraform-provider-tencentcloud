package ccn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCcnCrossBorderFlowMonitorDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnCrossBorderFlowMonitorDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_cross_border_flow_monitor.cross_border_flow_monitor")),
			},
		},
	})
}

const testAccCcnCrossBorderFlowMonitorDataSource = `

data "tencentcloud_ccn_cross_border_flow_monitor" "cross_border_flow_monitor" {
  source_region = "ap-guangzhou"
  destination_region = "ap-singapore"
  ccn_id = "ccn-39lqkygf"
  ccn_uin = "979137"
  period = 60
  start_time = "2023-01-01 00:00:00"
  end_time = "2023-01-01 01:00:00"
}

`
