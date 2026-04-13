package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTeoExportZoneConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoExportZoneConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoExportZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.export_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_config", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.export_config", "export_type", "all"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_config", "config_content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_export_zone_config.export_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoExportZoneConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoExportZoneConfigExists("tencentcloud_teo_export_zone_config.export_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_config", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_export_zone_config.export_config", "export_type", "basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_export_zone_config.export_config", "config_content"),
				),
			},
		},
	})
}

func testAccCheckTeoExportZoneConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_export_zone_config" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]

		// For export resource, we just verify the zone exists
		zone, err := service.DescribeTeoZone(ctx, zoneId)
		if zone != nil {
			// Zone still exists, which is expected for export resource
			return nil
		}
		if err != nil {
			return err
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

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		zone, err := service.DescribeTeoZone(ctx, zoneId)
		if zone == nil {
			return fmt.Errorf("Zone %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoExportZoneConfig = testAccTeoZone + `

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [tencentcloud_teo_zone.basic]
}

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
    zone_id     = tencentcloud_teo_zone.basic.id
    domain_name = var.zone_name

    origin_info {
        origin      = "150.109.8.1"
        origin_type = "IP_DOMAIN"
    }

    status            = "online"
    origin_protocol   = "FOLLOW"
    http_origin_port  = 80
    https_origin_port = 443
    ipv6_status       = "follow"

    depends_on = [tencentcloud_teo_ownership_verify.ownership_verify]
}

resource "tencentcloud_teo_export_zone_config" "export_config" {
    zone_id     = tencentcloud_teo_zone.basic.id
    export_type = "all"

    depends_on = [tencentcloud_teo_acceleration_domain.acceleration_domain]
}
`

const testAccTeoExportZoneConfigBasic = testAccTeoZone + `

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [tencentcloud_teo_zone.basic]
}

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
    zone_id     = tencentcloud_teo_zone.basic.id
    domain_name = var.zone_name

    origin_info {
        origin      = "150.109.8.1"
        origin_type = "IP_DOMAIN"
    }

    status            = "online"
    origin_protocol   = "FOLLOW"
    http_origin_port  = 80
    https_origin_port = 443
    ipv6_status       = "follow"

    depends_on = [tencentcloud_teo_ownership_verify.ownership_verify]
}

resource "tencentcloud_teo_export_zone_config" "export_config" {
    zone_id     = tencentcloud_teo_zone.basic.id
    export_type = "basic"

    depends_on = [tencentcloud_teo_acceleration_domain.acceleration_domain]
}
`
