package dnspod_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDnspodPackageOrderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodPackageOrder,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dnspod_package_order.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_package_order.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodPackageOrder = `
resource "tencentcloud_dnspod_package_order" "example" {
  domain = "demo.com"
  grade  = "DPG_ULTIMATE"
}
`
