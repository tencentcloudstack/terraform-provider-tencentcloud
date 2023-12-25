package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcSubnetResourceDashboardDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetResourceDashboardDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_subnet_resource_dashboard.subnet_resource_dashboard")),
			},
		},
	})
}

const testAccVpcSubnetResourceDashboardDataSource = `

data "tencentcloud_vpc_subnet_resource_dashboard" "subnet_resource_dashboard" {
  subnet_ids = ["subnet-i9tpf6hq"]
}

`
