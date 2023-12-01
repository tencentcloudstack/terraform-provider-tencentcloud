package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcNetworkInterfaceLimitDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetworkInterfaceLimitDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_network_interface_limit.network_interface_limit")),
			},
		},
	})
}

const testAccVpcNetworkInterfaceLimitDataSource = `

data "tencentcloud_vpc_network_interface_limit" "network_interface_limit" {
  instance_id = "ins-cr2rfq78"
}

`
