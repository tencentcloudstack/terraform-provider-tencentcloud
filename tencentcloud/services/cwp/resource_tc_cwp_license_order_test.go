package cwp_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCwpLicenseOrderResource_basic -v
func TestAccTencentCloudNeedFixCwpLicenseOrderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCwpLicenseOrder,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "alias"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "license_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "license_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "region_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "project_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cwp_license_order.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCwpLicenseOrderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "alias"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "license_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "license_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "region_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_order.example", "project_id"),
				),
			},
		},
	})
}

const testAccCwpLicenseOrder = `
resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags        = {
    "createdBy" = "terraform"
  }
}
`

const testAccCwpLicenseOrderUpdate = `
resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example1"
  license_type = 0
  license_num  = 2
  region_id    = 1
  project_id   = 0
  tags        = {
    "createdBy1" = "terraform1"
  }
}
`
