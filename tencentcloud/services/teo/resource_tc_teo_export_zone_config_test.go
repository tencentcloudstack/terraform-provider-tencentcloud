package teo_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -v ./tencentcloud/services/teo -run TestAccTencentcloudTeoExportZoneConfig_basic
func TestAccTencentcloudTeoExportZoneConfig_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.basic", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.basic", "content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_export_zone_config.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// go test -v ./tencentcloud/services/teo -run TestAccTencentcloudTeoExportZoneConfig_withTypes
func TestAccTencentcloudTeoExportZoneConfig_withTypes(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigWithTypes,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.with_types"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.with_types", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.with_types", "types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.with_types", "types.0", "L7AccelerationConfig"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.with_types", "content"),
				),
			},
		},
	})
}

// go test -v ./tencentcloud/services/teo -run TestAccTencentcloudTeoExportZoneConfig_update
func TestAccTencentcloudTeoExportZoneConfig_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.basic", "types.#", "0"),
				),
			},
			{
				Config:      testAccTeoExportZoneConfigWithTypes,
				ExpectError: regexp.MustCompile(`does not support update, changes require recreation`),
			},
		},
	})
}

// go test -v ./tencentcloud/services/teo -run TestAccTencentcloudTeoExportZoneConfig_delete
func TestAccTencentcloudTeoExportZoneConfig_delete(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.basic"),
				),
			},
			{
				Config: testAccTeoExportZoneConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					// After deletion, resource should be gone from state
				),
			},
		},
	})
}

// go test -v ./tencentcloud/services/teo -run TestAccTencentcloudTeoExportZoneConfig_timeout
func TestAccTencentcloudTeoExportZoneConfig_timeout(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigWithTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.timeout"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.timeout", "zone_id"),
				),
			},
		},
	})
}

// go test -v ./tencentcloud/services/teo -run TestAccTencentcloudTeoExportZoneConfig_importState
func TestAccTencentcloudTeoExportZoneConfig_importState(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.basic"),
				),
			},
			{
				ResourceName:            "tencentcloud_teo_export_zone_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"types", "content"},
			},
		},
	})
}

func testAccCheckTeoExportZoneConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	// ExportZoneConfig is a read-only resource, logical deletion only
	// This check verifies that the zone still exists (the export itself is not a persistent resource)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_export_zone_config" {
			continue
		}

		zoneId := rs.Primary.Attributes["zone_id"]
		zone, err := service.DescribeTeoZone(ctx, zoneId)
		if err != nil {
			return fmt.Errorf("error checking zone %s existence: %s", zoneId, err)
		}
		if zone == nil {
			return fmt.Errorf("zone %s does not exist", zoneId)
		}
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

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		zoneId := rs.Primary.Attributes["zone_id"]
		zone, err := service.DescribeTeoZone(ctx, zoneId)
		if zone == nil {
			return fmt.Errorf("zone %s is not found", zoneId)
		}
		if err != nil {
			return err
		}

		// Verify content is set
		if _, ok := rs.Primary.Attributes["content"]; !ok {
			return fmt.Errorf("export zone config content is empty")
		}

		return nil
	}
}

const testAccTeoExportZoneConfigVar = `
variable "zone_id" {
  description = "Zone ID for testing"
  default = ""
}
`

const testAccTeoExportZoneConfigBasic = testAccTeoExportZoneConfigVar + `
resource "tencentcloud_teo_export_zone_config" "basic" {
  zone_id = var.zone_id
}
`

const testAccTeoExportZoneConfigWithTypes = testAccTeoExportZoneConfigVar + `
resource "tencentcloud_teo_export_zone_config" "with_types" {
  zone_id = var.zone_id
  types   = ["L7AccelerationConfig"]
}
`

const testAccTeoExportZoneConfigWithTimeout = testAccTeoExportZoneConfigVar + `
resource "tencentcloud_teo_export_zone_config" "timeout" {
  zone_id = var.zone_id

  timeouts {
    create = "10m"
    read   = "10m"
  }
}
`

const testAccTeoExportZoneConfigEmpty = testAccTeoExportZoneConfigVar + `
# Empty config - resource is deleted
`
