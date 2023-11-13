package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcDcxExtraConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcDcxExtraConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dc_dcx_extra_config.dcx_extra_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dc_dcx_extra_config.dcx_extra_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcDcxExtraConfig = `

resource "tencentcloud_dc_dcx_extra_config" "dcx_extra_config" {
  direct_connect_tunnel_id = "dcx-test123"
  vlan = 123
  bgp_peer {
		asn = 65101
		auth_key = "test123"

  }
  route_filter_prefixes {
		cidr = "192.168.0.0/24"

  }
  tencent_address = "192.168.1.1"
  tencent_backup_address = "192.168.1.2"
  customer_address = "192.168.1.4"
  bandwidth = 10M
  enable_b_g_p_community = false
  bfd_enable = false
  nqa_enable = false
  bfd_info {
		probe_failed_times = 3
		interval = 100

  }
  nqa_info {
		probe_failed_times = 3
		interval = 100
		destination_ip = "192.168.2.2"

  }
  i_pv6_enable = 0
  customer_i_d_c_routes {
		cidr = ""

  }
  jumbo_enable = 0
}

`
