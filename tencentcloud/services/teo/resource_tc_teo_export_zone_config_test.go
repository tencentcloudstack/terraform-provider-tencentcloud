package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoExportZoneConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "content"),
				),
			},
			{
				Config: testAccTeoExportZoneConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_export_zone_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoExportZoneConfig = `
resource "tencentcloud_teo_zone" "example" {
  zone_name = "terraform-test.example.com"
  type      = "generic"
}

resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = tencentcloud_teo_zone.example.id
}
`

const testAccTeoExportZoneConfigUpdate = `
resource "tencentcloud_teo_zone" "example" {
  zone_name = "terraform-test.example.com"
  type      = "generic"
}

resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = tencentcloud_teo_zone.example.id
  types   = ["L7AccelerationConfig"]
}
`
