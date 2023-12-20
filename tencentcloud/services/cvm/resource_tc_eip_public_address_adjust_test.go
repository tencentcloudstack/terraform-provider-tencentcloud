package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixVpcPublicAddressAdjustResource_basic -v
func TestAccTencentCloudNeedFixVpcPublicAddressAdjustResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPublicAddressAdjust,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_eip_public_address_adjust.public_address_adjust", "id"),
				),
			},
		},
	})
}

const testAccVpcPublicAddressAdjust = `
resource "tencentcloud_eip_public_address_adjust" "public_address_adjust" {
  instance_id = "ins-cr2rfq78"
  address_id  = "eip-erft45fu"
}
`
