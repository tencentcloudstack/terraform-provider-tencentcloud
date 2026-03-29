package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoExportZoneConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoExportZoneConfigDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_export_zone_config.export_config"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_export_zone_config.export_config", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.export_config", "zone_name"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.export_config", "area"),
			),
		}},
	})
}

func TestAccTencentCloudTeoExportZoneConfigDataSource_invalidZoneId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config:      testAccTeoExportZoneConfigDataSourceInvalidZoneId,
			ExpectError: tcacctest.RegexpMatches(`.*zone setting not found.*`),
		}},
	})
}

const testAccTeoExportZoneConfigDataSource = `

data "tencentcloud_teo_export_zone_config" "export_config" {
  zone_id = "zone-2xkazzl8yf6k"
}
`

const testAccTeoExportZoneConfigDataSourceInvalidZoneId = `

data "tencentcloud_teo_export_zone_config" "export_config" {
  zone_id = "zone-nonexistent"
}
`
