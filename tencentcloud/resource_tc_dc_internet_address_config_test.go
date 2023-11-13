package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcInternetAddressConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcInternetAddressConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dc_internet_address_config.internet_address_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dc_internet_address_config.internet_address_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcInternetAddressConfig = `

resource "tencentcloud_dc_internet_address_config" "internet_address_config" {
  instance_id = "ipv4-ljm17pbl"
  enable = true
}

`
