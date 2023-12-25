package vpn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnDefaultHealthCheckIpDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnDefaultHealthCheckIpDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_default_health_check_ip.default_health_check_ip")),
			},
		},
	})
}

const testAccVpnDefaultHealthCheckIpDataSource = `

data "tencentcloud_vpn_default_health_check_ip" "default_health_check_ip" {
  vpn_gateway_id = "vpngw-gt8bianl"
}

`
