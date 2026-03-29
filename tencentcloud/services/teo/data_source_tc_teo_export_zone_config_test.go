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
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_export_zone_config.config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "zone_id_output"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "zone_name_output"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "area"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "status"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoExportZoneConfigDataSource_byZoneId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigDataSourceByZoneId,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_export_zone_config.config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "zone_id_output"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoExportZoneConfigDataSource_byZoneName(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigDataSourceByZoneName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_export_zone_config.config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_teo_export_zone_config.config", "zone_name_output"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoExportZoneConfigDataSource_zoneIdPriority(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigDataSourceZoneIdPriority,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_export_zone_config.config"),
				),
			},
		},
	})
}

const testAccTeoExportZoneConfigDataSource = `
data "tencentcloud_teo_export_zone_config" "config" {
  zone_id = "zone-xxxxx"
}
`

const testAccTeoExportZoneConfigDataSourceByZoneId = `
data "tencentcloud_teo_export_zone_config" "config" {
  zone_id = "zone-xxxxx"
}
`

const testAccTeoExportZoneConfigDataSourceByZoneName = `
data "tencentcloud_teo_export_zone_config" "config" {
  zone_name = "example.com"
}
`

const testAccTeoExportZoneConfigDataSourceZoneIdPriority = `
data "tencentcloud_teo_export_zone_config" "config" {
  zone_id   = "zone-xxxxx"
  zone_name = "example.com"
}
`
