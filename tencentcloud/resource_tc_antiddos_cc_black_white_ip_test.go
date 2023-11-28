package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosCcBlackWhiteIpResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosCcBlackWhiteIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "black_white_ip.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "domain", "t.baidu.com"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "instance_id", "bgpip-0000078h"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "ip", "212.64.62.191"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip", "type", "black"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosCcBlackWhiteIp = `

resource "tencentcloud_antiddos_cc_black_white_ip" "cc_black_white_ip" {
	instance_id = "bgpip-0000078h"
	black_white_ip {
	  ip   = "1.2.3.5"
	  mask = 0
  
	}
	type     = "black"
	ip       = "212.64.62.191"
	domain   = "t.baidu.com"
	protocol = "http"
  }

`
