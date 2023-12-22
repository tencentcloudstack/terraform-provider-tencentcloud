package vpn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnCustomerGatewayVendorsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayVendorsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_customer_gateway_vendors.vpn_customer_gateway_vendors")),
			},
		},
	})
}

const testAccVpnCustomerGatewayVendorsDataSource = `

data "tencentcloud_vpn_customer_gateway_vendors" "vpn_customer_gateway_vendors" {}

`
