package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixEipAddressTransformResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAddressTransform,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eip_address_transform.address_transform", "id")),
			},
			{
				ResourceName:      "tencentcloud_eip_address_transform.address_transform",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEipAddressTransform = `

resource "tencentcloud_eip_address_transform" "address_transform" {
  instance_id = "ins-2kcdugsq"
}

`
