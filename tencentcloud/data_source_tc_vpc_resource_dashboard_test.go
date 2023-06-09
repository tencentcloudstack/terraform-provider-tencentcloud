package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcResourceDashboardDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcResourceDashboardDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_resource_dashboard.resource_dashboard")),
			},
		},
	})
}

const testAccVpcResourceDashboardDataSource = `

data "tencentcloud_vpc_resource_dashboard" "resource_dashboard" {
  vpc_ids = ["vpc-4owdpnwr"]
}

`
