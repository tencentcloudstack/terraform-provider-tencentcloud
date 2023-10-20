package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapCheckProxyCreateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCheckProxyCreateDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_check_proxy_create.check_proxy_create"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_check_proxy_create.check_proxy_create", "check_flag", "1"),
				),
			},
		},
	})
}

const testAccGaapCheckProxyCreateDataSource = `
data "tencentcloud_gaap_check_proxy_create" "check_proxy_create" {
  access_region = "Guangzhou"
  real_server_region = "Beijing"
  bandwidth = 10
  concurrent = 2
  ip_address_version = "IPv4"
  network_type = "normal"
  package_type = "Thunder"
}
`
