package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy_group.proxy_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy_group.proxy_group", "group_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy_group.proxy_group", "ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy_group.proxy_group", "package_type", "Thunder"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy_group.proxy_group", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy_group.proxy_group", "real_server_region", "Beijing"),
				),
			},
			{
				Config: testAccGaapProxyGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy_group.proxy_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy_group.proxy_group", "group_name", "tf-test-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_proxy_group.proxy_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGaapProxyGroup = `
resource "tencentcloud_gaap_proxy_group" "proxy_group" {
  project_id = 0
  group_name = "tf-test"
  real_server_region = "Beijing"
  ip_address_version = "IPv4"
  package_type = "Thunder"
}
`

const testAccGaapProxyGroupUpdate = `
resource "tencentcloud_gaap_proxy_group" "proxy_group" {
  project_id = 0
  group_name = "tf-test-update"
  real_server_region = "Beijing"
  ip_address_version = "IPv4"
  package_type = "Thunder"
}
`
