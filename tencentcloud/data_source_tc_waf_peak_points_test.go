package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafPeakPointsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPeakPointsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_peak_points.peak_points")),
			},
		},
	})
}

const testAccWafPeakPointsDataSource = `

data "tencentcloud_waf_peak_points" "peak_points" {
  from_time = ""
  to_time = ""
  domain = ""
  edition = ""
  instance_i_d = ""
  metric_name = ""
  }

`
