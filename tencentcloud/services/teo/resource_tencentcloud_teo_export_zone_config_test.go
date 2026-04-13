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

func TestAccTencentCloudTeoExportZoneConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.export_zone_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_zone_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_zone_config", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_zone_config", "content"),
				),
			},
			{
				Config: testAccTeoExportZoneConfigWithTypes,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.export_zone_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_zone_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_zone_config", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.export_zone_config", "types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.export_zone_config", "types.0", "L7AccelerationConfig"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_zone_config", "content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_export_zone_config.export_zone_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTeoExportZoneConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_export_zone_config" {
			continue
		}

		// This is an export resource, so we don't need to check if it exists in the API
		// Just verify it's removed from Terraform state
		log.Printf("[DEBUG]%s export zone config destroyed, zone_id [%s]\n", logId, rs.Primary.ID)
	}

	return nil
}

func testAccCheckTeoExportZoneConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		// This is an export resource, we just verify it has a valid ID and content
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID", r)
		}

		if _, ok := rs.Primary.Attributes["zone_id"]; !ok {
			return fmt.Errorf("resource %s has no zone_id", r)
		}

		if _, ok := rs.Primary.Attributes["content"]; !ok {
			return fmt.Errorf("resource %s has no content", r)
		}

		return nil
	}
}

const testAccTeoExportZoneConfigBasic = testAccTeoZone + `

resource "tencentcloud_teo_export_zone_config" "export_zone_config" {
  zone_id = tencentcloud_teo_zone.basic.id
}
`

const testAccTeoExportZoneConfigWithTypes = testAccTeoZone + `

resource "tencentcloud_teo_export_zone_config" "export_zone_config" {
  zone_id = tencentcloud_teo_zone.basic.id

  types = [
    "L7AccelerationConfig",
  ]
}
`
