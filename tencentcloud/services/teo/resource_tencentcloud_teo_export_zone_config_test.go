package teo_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -test.run TestAccTencentCloudTeoExportZoneConfig_basic -v
func TestAccTencentCloudTeoExportZoneConfig_basic(t *testing.T) {
	t.Parallel()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.PreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.test", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.test", "content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_export_zone_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// go test -test.run TestAccTencentCloudTeoExportZoneConfig_withTypes -v
func TestAccTencentCloudTeoExportZoneConfig_withTypes(t *testing.T) {
	t.Parallel()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.PreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigWithTypes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.test", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.test", "content"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.test", "types.0", "L7AccelerationConfig"),
				),
			},
		},
	})
}

func testAccCheckTeoExportZoneConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cli, _ := tcacctest.SharedClientForRegion(tcacctest.DefaultRegion)
		client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcteo.NewTeoService(client)

		zoneId := rs.Primary.ID
		types := make([]*string, 0)
		if v, ok := rs.Primary.Attributes["types.0"]; ok && v != "" {
			types = append(types, &v)
		}

		respData, err := service.ExportTeoZoneConfigById(ctx, zoneId, types)
		if err != nil {
			return err
		}

		if respData == nil || respData.Content == nil {
			return fmt.Errorf("export zone config response is empty")
		}

		return nil
	}
}

func testAccTeoExportZoneConfigBasic() string {
	return `
resource "tencentcloud_teo_zone" "test" {
  type     = "partial"
  area     = "overseas"
  plan_id  = "plan-00000000"
  zone_name = "test.example.com"
}

resource "tencentcloud_teo_export_zone_config" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  depends_on = [tencentcloud_teo_zone.test]
}
`
}

func testAccTeoExportZoneConfigWithTypes() string {
	return `
resource "tencentcloud_teo_zone" "test" {
  type     = "partial"
  area     = "overseas"
  plan_id  = "plan-00000000"
  zone_name = "test2.example.com"
}

resource "tencentcloud_teo_export_zone_config" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  types   = ["L7AccelerationConfig"]
  depends_on = [tencentcloud_teo_zone.test]
}
`
}
