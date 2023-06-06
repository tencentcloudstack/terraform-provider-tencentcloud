package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcUsedIpAddressDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcUsedIpAddressDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_used_ip_address.used_ip_address")),
			},
		},
	})
}

const testAccVpcUsedIpAddressDataSource = `

data "tencentcloud_vpc_used_ip_address" "used_ip_address" {
  vpc_id = "vpc-4owdpnwr"
}

`
