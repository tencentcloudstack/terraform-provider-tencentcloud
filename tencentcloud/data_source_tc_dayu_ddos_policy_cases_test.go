package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataDayuDdosPolicyCasesName = "data.tencentcloud_dayu_ddos_policy_cases.id_test"

func TestAccTencentCloudDataDayuDdosPolicyCases(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyCaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataDayuDdosPolicyCasesBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuDdosPolicyCaseExists("tencentcloud_dayu_ddos_policy_case.test_policy_case"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.name", "tf_test_policy_case"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.platform_types.#", "2"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.web_api_urls.#", "2"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.web_api_urls.0", "abc.com"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.app_protocols.#", "2"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.has_abroad", "yes"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.has_initiate_tcp", "yes"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.has_initiate_udp", "yes"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.has_vpn", "yes"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.tcp_start_port", "1000"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.tcp_end_port", "2000"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.udp_start_port", "3000"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.udp_end_port", "4000"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.peer_tcp_port", "1111"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.peer_udp_port", "3333"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.tcp_footprint", "511"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.udp_footprint", "500"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.min_tcp_package_len", "1000"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.max_tcp_package_len", "1200"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.min_udp_package_len", "1000"),
					resource.TestCheckResourceAttr(testDataDayuDdosPolicyCasesName, "list.0.max_udp_package_len", "1200"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuDdosPolicyCasesBaic = `
resource "tencentcloud_dayu_ddos_policy_case" "test_policy_case" {
  resource_type         = "bgpip"
  name                  = "tf_test_policy_case"
  platform_types        = ["PC", "MOBILE"]
  app_type              = "WEB"
  app_protocols         = ["tcp", "udp"]
  tcp_start_port		= "1000"
  tcp_end_port          = "2000"
  udp_start_port		= "3000"
  udp_end_port			= "4000"
  has_abroad			= "yes"
  has_initiate_tcp		= "yes"
  has_initiate_udp		= "yes"
  peer_tcp_port			= "1111"
  peer_udp_port			= "3333"
  tcp_footprint		= "511"
  udp_footprint		= "500"
  web_api_urls			= ["abc.com", "test.cn/aaa.png"]
  min_tcp_package_len	= "1000"
  max_tcp_package_len	= "1200"
  min_udp_package_len	= "1000"
  max_udp_package_len	= "1200"
  has_vpn				= "yes"
}

data "tencentcloud_dayu_ddos_policy_cases" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy_case.test_policy_case.resource_type
  scene_id      = tencentcloud_dayu_ddos_policy_case.test_policy_case.scene_id
}
`
