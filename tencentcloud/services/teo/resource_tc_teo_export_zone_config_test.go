package teo_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoExportZoneConfig_basic(t *testing.T) {
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
				Config: testAccTeoExportZoneConfigWithTypes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.example", "content"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.example", "export_types.#", "1"),
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

func TestAccTencentCloudTeoExportZoneConfig_disappears(t *testing.T) {
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
				),
				ExpectNonEmptyPlan: true, // Since this is a virtual resource, it will always show up in the plan
			},
		},
	})
}

var testZoneId string

func init() {
	// This should be set to a real zone ID for testing
	// Example: "zone-xxxxxxxxxxxxx"
	testZoneId = "zone-xxxxxxxxxxxxx"
}

const testAccTeoExportZoneConfig = fmt.Sprintf(`
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "%s"
}
`, testZoneId)

const testAccTeoExportZoneConfigWithTypes = fmt.Sprintf(`
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "%s"
  export_types = ["L7AccelerationConfig"]
}
`, testZoneId)
