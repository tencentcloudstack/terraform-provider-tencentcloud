package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorPromRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorPromRegionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_prom_regions.prom_regions")),
			},
		},
	})
}

const testAccMonitorPromRegionsDataSource = `

data "tencentcloud_monitor_prom_regions" "prom_regions" {
  pay_mode = 
  }

`
