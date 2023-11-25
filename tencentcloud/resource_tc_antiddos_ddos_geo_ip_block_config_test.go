package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosDdosGeoIpBlockConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosDdosGeoIpBlockConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "ddos_geo_ip_block_config.0.action", "drop"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "ddos_geo_ip_block_config.0.region_type", "customized"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "ddos_geo_ip_block_config.0.area_list.#", "1"),
				),
			},
			{
				Config: testAccAntiddosDdosGeoIpBlockConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "ddos_geo_ip_block_config.0.action", "drop"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "ddos_geo_ip_block_config.0.region_type", "customized"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config", "ddos_geo_ip_block_config.0.area_list.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_ddos_geo_ip_block_config.ddos_geo_ip_block_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosDdosGeoIpBlockConfig = `
resource "tencentcloud_antiddos_ddos_geo_ip_block_config" "ddos_geo_ip_block_config" {
	instance_id = "bgp-00000ry7"
	ddos_geo_ip_block_config {
	  region_type = "customized"
	  action = "drop"
	  area_list = [100002]
	}
}
`

const testAccAntiddosDdosGeoIpBlockConfigUpdate = `
resource "tencentcloud_antiddos_ddos_geo_ip_block_config" "ddos_geo_ip_block_config" {
	instance_id = "bgp-00000ry7"
	ddos_geo_ip_block_config {
	  region_type = "customized"
	  action = "drop"
	  area_list = [100002, 100003]
	}
}
`
