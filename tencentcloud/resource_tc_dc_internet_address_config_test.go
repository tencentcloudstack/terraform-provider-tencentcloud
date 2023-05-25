package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixDcInternetAddressConfigResource_basic(t *testing.T) {
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

resource "tencentcloud_dc_internet_address" "internet_address" {
  mask_len = 30
  addr_type = 2
  addr_proto = 0
}

resource "tencentcloud_dc_internet_address_config" "internet_address_config" {
  instance_id = tencentcloud_dc_internet_address.internet_address.id
  enable = false
}

`
