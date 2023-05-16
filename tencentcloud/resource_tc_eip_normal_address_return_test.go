package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixEipNormalAddressReturnResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipNormalAddressReturn,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eip_normal_address_return.normal_address_return", "id")),
			},
		},
	})
}

const testAccEipNormalAddressReturn = `

resource "tencentcloud_eip_normal_address_return" "normal_address_return" {
  address_ips = ["111.230.44.68"]
}
`
