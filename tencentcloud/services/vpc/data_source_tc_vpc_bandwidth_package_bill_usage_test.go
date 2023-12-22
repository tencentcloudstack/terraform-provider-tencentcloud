package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcBandwidthPackageBillUsageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackageBillUsageDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_bandwidth_package_bill_usage.bandwidth_package_bill_usage")),
			},
		},
	})
}

const testAccVpcBandwidthPackageBillUsageDataSource = `

resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type            = "BGP"
  charge_type             = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name  = "iac-test-data"
  tags = {
    "createdBy" = "terraform"
  }
}

data "tencentcloud_vpc_bandwidth_package_bill_usage" "bandwidth_package_bill_usage" {
  bandwidth_package_id =  tencentcloud_vpc_bandwidth_package.bandwidth_package.id
}

`
