package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcVpnCustomerGatewayVendorsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnCustomerGatewayVendorsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_vpn_customer_gateway_vendors.vpn_customer_gateway_vendors")),
			},
		},
	})
}

const testAccVpcVpnCustomerGatewayVendorsDataSource = `

data "tencentcloud_vpc_vpn_customer_gateway_vendors" "vpn_customer_gateway_vendors" {
  }

`
