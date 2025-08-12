package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoDdosProtectionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDdosProtectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_ddos_protection_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_ddos_protection_config.example", "zone_id"),
				),
			},
			{
				Config: testAccTeoDdosProtectionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_ddos_protection_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_ddos_protection_config.example", "zone_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_ddos_protection_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoDdosProtectionConfig = `
resource "tencentcloud_teo_ddos_protection_config" "example" {
  zone_id = "zone-3edjdliiw3he"
  ddos_protection {
    protection_option = "protect_all_domains"
  }
}
`

const testAccTeoDdosProtectionConfigUpdate = `
resource "tencentcloud_teo_ddos_protection_config" "example" {
  zone_id = "zone-3edjdliiw3he"
  ddos_protection {
    protection_option = "protect_specified_domains"
    domain_ddos_protections {
      domain = "1.demo.com"
      switch = "on"
    }

    domain_ddos_protections {
      domain = "2.demo.com"
      switch = "on"
    }

    domain_ddos_protections {
      domain = "3.demo.com"
      switch = "off"
    }
  }
}
`
