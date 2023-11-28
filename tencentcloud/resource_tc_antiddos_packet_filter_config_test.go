package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosPacketFilterConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosPacketFilterConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_packet_filter_config.packet_filter_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_packet_filter_config.packet_filter_config", "instance_id", "bgp-00000ry7"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_packet_filter_config.packet_filter_config", "packet_filter_config.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_packet_filter_config.packet_filter_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosPacketFilterConfig = `
resource "tencentcloud_antiddos_packet_filter_config" "packet_filter_config" {
	instance_id = "bgp-00000ry7"
	packet_filter_config {
	  action      = "drop"
	  depth       = 1
	  dport_start = 80
	  dport_end   = 80
	  is_not      = 0
	  match_begin = "begin_l5"
	  match_type  = "pcre"
	  offset      = 1
	  pktlen_min  = 1400
	  pktlen_max  = 1400
	  protocol    = "all"
	  sport_start = 8080
	  sport_end   = 8080
	  str         = "a"
	}
  }
`
