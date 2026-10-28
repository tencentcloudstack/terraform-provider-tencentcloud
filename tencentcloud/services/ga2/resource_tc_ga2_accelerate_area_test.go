package ga2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGa2AccelerateAreaResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGa2AccelerateArea,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_accelerate_area.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_accelerate_area.example", "accelerator_area_id"),
				),
			},
			{
				Config: testAccGa2AccelerateAreaUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ga2_accelerate_area.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ga2_accelerate_area.example", "bandwidth", "20"),
				),
			},
			{
				ResourceName:      "tencentcloud_ga2_accelerate_area.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGa2AccelerateArea = `
resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = "ga-4mredmiu"
  accelerate_region     = "ap-guangzhou"
  bandwidth             = 10
  isp_type              = "BGP"
  ip_version            = "IPv4"
}
`

const testAccGa2AccelerateAreaUpdate = `
resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = "ga-4mredmiu"
  accelerate_region     = "ap-guangzhou"
  bandwidth             = 20
  isp_type              = "BGP"
  ip_version            = "IPv4"
}
`
