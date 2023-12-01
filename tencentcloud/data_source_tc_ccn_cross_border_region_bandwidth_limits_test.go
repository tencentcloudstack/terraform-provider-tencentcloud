package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCcnCrossBorderRegionBandwidthLimitsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnCrossBorderRegionBandwidthLimitsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_cross_border_region_bandwidth_limits.ccn_region_bandwidth_limits")),
			},
		},
	})
}

const testAccCcnCrossBorderRegionBandwidthLimitsDataSource = `

data "tencentcloud_ccn_cross_border_region_bandwidth_limits" "ccn_region_bandwidth_limits" {
  filters {
    name   = "source-region"
    values = ["ap-guangzhou"]
  }

  filters {
    name   = "destination-region"
    values = ["ap-shanghai"]
  }
}

`
