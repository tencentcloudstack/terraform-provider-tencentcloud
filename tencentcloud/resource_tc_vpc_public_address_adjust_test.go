package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcPublicAddressAdjustResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPublicAddressAdjust,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_public_address_adjust.public_address_adjust", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_public_address_adjust.public_address_adjust",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPublicAddressAdjust = `

resource "tencentcloud_vpc_public_address_adjust" "public_address_adjust" {
  instance_id = "ins-osckfnm7"
  address_id = "eip-erft45fu"
}

`
