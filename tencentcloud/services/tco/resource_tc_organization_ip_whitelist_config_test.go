package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationIPWhitelistConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationIPWhitelistConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_ip_whitelist_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_ip_whitelist_config.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_ip_whitelist_config.example", "ip_whitelist"),
				),
			},
			{
				Config: testAccOrganizationIPWhitelistConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_ip_whitelist_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_ip_whitelist_config.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_ip_whitelist_config.example", "ip_whitelist"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_ip_whitelist_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationIPWhitelistConfig = `
resource "tencentcloud_organization_ip_whitelist_config" "example" {
  zone_id = "z-1os7c9znogct"
  ip_whitelist = [
    "10.0.0.0/24",
    "192.168.1.0/24",
    "172.16.10.0/24",
  ]
}
`

const testAccOrganizationIPWhitelistConfigUpdate = `
resource "tencentcloud_organization_ip_whitelist_config" "example" {
  zone_id = "z-1os7c9znogct"
  ip_whitelist = [
    "10.0.0.0/24",
  ]
}
`
