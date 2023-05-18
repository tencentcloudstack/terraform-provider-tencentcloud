package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcTrafficPackageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcTrafficPackage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_traffic_package.traffic_package", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_traffic_package.traffic_package",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcTrafficPackage = `

resource "tencentcloud_vpc_traffic_package" "traffic_package" {
  traffic_amount = 10
}
`
