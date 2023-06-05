package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcPrivateIpAddressesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPrivateIpAddressesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_private_ip_addresses.private_ip_addresses")),
			},
		},
	})
}

const testAccVpcPrivateIpAddressesDataSource = `

data "tencentcloud_vpc_private_ip_addresses" "private_ip_addresses" {
  vpc_id = "vpc-l0dw94uh"
  private_ip_addresses = ["10.0.0.1"]
}

`
