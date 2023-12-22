package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosDdosBlackWhiteIpResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosDdosBlackWhiteIp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip", "mask", "0"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip", "type", "black"),
				),
			},
			{
				Config: testAccAntiddosDdosBlackWhiteIpUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip", "mask", "0"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip", "type", "white"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosDdosBlackWhiteIp = `
resource "tencentcloud_antiddos_ddos_black_white_ip" "ddos_black_white_ip" {
	instance_id = "bgp-00000ry7"
	ip = "1.2.3.5"
	mask = 0
	type = "black"
}
`

const testAccAntiddosDdosBlackWhiteIpUpdate = `
resource "tencentcloud_antiddos_ddos_black_white_ip" "ddos_black_white_ip" {
	instance_id = "bgp-00000ry7"
	ip = "1.2.3.5"
	mask = 0
	type = "white"
}
`
