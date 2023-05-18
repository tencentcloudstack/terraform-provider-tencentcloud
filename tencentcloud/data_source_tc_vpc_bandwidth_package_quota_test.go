package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcBandwidthPackageQuotaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackageQuotaDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_bandwidth_package_quota.bandwidth_package_quota")),
			},
		},
	})
}

const testAccVpcBandwidthPackageQuotaDataSource = `

data "tencentcloud_vpc_bandwidth_package_quota" "bandwidth_package_quota" {
  }

`
