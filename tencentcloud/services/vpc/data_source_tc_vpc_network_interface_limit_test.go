package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcNetworkInterfaceLimitDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetworkInterfaceLimitDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_network_interface_limit.network_interface_limit")),
			},
		},
	})
}

const testAccVpcNetworkInterfaceLimitDataSource = `

data "tencentcloud_vpc_network_interface_limit" "network_interface_limit" {
  instance_id = "ins-cr2rfq78"
}

`
