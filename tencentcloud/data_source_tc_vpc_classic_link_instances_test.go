package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcClassicLinkInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcClassicLinkInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_classic_link_instances.classic_link_instances")),
			},
		},
	})
}

const testAccVpcClassicLinkInstancesDataSource = `

data "tencentcloud_vpc_classic_link_instances" "classic_link_instances" {
  filters {
    name   = "vpc-id"
    values = ["vpc-lh4nqig9"]
  }
}

`
