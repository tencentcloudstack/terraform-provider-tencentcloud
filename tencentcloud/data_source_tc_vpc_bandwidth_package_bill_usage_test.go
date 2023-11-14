package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcBandwidthPackageBillUsageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackageBillUsageDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_bandwidth_package_bill_usage.bandwidth_package_bill_usage")),
			},
		},
	})
}

const testAccVpcBandwidthPackageBillUsageDataSource = `

data "tencentcloud_vpc_bandwidth_package_bill_usage" "bandwidth_package_bill_usage" {
  bandwidth_package_id = "bwp-234rfgt5"
  }

`
