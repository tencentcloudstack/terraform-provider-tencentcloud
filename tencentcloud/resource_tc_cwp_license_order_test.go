package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCwpLicenseOrderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCwpLicenseOrder,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.license_order", "id")),
			},
			{
				ResourceName:      "tencentcloud_cwp_license_order.license_order",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCwpLicenseOrder = `

resource "tencentcloud_cwp_license_order" "license_order" {
  license_type = 
  license_num = 
  region_id = 
  project_id = 
}

`
