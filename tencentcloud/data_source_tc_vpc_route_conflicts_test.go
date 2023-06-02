package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcRouteConflictsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteConflictsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_route_conflicts.route_conflicts")),
			},
		},
	})
}

const testAccVpcRouteConflictsDataSource = `

data "tencentcloud_vpc_route_conflicts" "route_conflicts" {
  route_table_id = "rtb-6xypllqe"
  destination_cidr_blocks = ["172.18.111.0/24"]
}

`
