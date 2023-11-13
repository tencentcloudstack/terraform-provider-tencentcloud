package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcAddressTransformResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcAddressTransform,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_address_transform.address_transform", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_address_transform.address_transform",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcAddressTransform = `

resource "tencentcloud_vpc_address_transform" "address_transform" {
  instance_id = ""
}

`
