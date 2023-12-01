package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnCustomerGatewayVendorsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayVendorsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_customer_gateway_vendors.vpn_customer_gateway_vendors")),
			},
		},
	})
}

const testAccVpnCustomerGatewayVendorsDataSource = `

data "tencentcloud_vpn_customer_gateway_vendors" "vpn_customer_gateway_vendors" {}

`
