package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcNormalAddressReturnResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNormalAddressReturn,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_normal_address_return.normal_address_return", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_normal_address_return.normal_address_return",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcNormalAddressReturn = `

resource "tencentcloud_vpc_normal_address_return" "normal_address_return" {
  address_ips = 
}

`
