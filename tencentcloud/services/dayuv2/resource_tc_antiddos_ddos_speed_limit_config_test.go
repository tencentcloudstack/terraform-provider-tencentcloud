package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosDdosSpeedLimitConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosDdosSpeedLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.protocol_list", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.dst_port_list", "8000"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.speed_values.0.type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.speed_values.0.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.speed_values.1.type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.speed_values.1.value", "2"),
				),
			},
			{
				Config: testAccAntiddosDdosSpeedLimitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.dst_port_list", "9000"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.speed_values.0.value", "0"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config", "ddos_speed_limit_config.0.speed_values.1.value", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosDdosSpeedLimitConfig = `
resource "tencentcloud_antiddos_ddos_speed_limit_config" "ddos_speed_limit_config" {
	instance_id = "bgp-00000ry7"
	ddos_speed_limit_config {
		  mode = 1
		  speed_values {
			  type = 1
			  value = 1
		  }
		  speed_values {
			  type = 2
			  value = 2
		  }
		  protocol_list = "ALL"
		  dst_port_list = "8000"
	}
}
`

const testAccAntiddosDdosSpeedLimitConfigUpdate = `
resource "tencentcloud_antiddos_ddos_speed_limit_config" "ddos_speed_limit_config" {
	instance_id = "bgp-00000ry7"
	ddos_speed_limit_config {
		  mode = 1
		  speed_values {
			  type = 1
			  value = 0
		  }
		  speed_values {
			  type = 2
			  value = 1
		  }
		  protocol_list = "ALL"
		  dst_port_list = "9000"
	}
}
`
