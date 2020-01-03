package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

var testDayuDdosPolicyCaseResourceName = "tencentcloud_dayu_ddos_policy_case"
var testDayuDdosPolicyCaseResourceKey = testDayuDdosPolicyCaseResourceName + ".test_policy_case"

func TestAccTencentCloudDayuDdosPolicyCaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyCaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuDdosPolicyCase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyCaseExists(testDayuDdosPolicyCaseResourceKey),
					resource.TestCheckResourceAttrSet(testDayuDdosPolicyCaseResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuDdosPolicyCaseResourceKey, "scene_id"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "resource_type", "bgp"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "name", "tf_test_policy_case"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "platform_types.#", "2"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "web_api_urls.#", "2"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "app_protocols.#", "2"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_abroad", "yes"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_initiate_tcp", "yes"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_initiate_udp", "yes"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_vpn", "yes"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "tcp_start_port", "1000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "tcp_end_port", "2000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "udp_start_port", "3000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "udp_end_port", "4000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "peer_tcp_port", "1111"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "peer_udp_port", "3333"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "tcp_footprint", "511"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "udp_footprint", "500"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "min_tcp_package_len", "1000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "max_tcp_package_len", "1200"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "min_udp_package_len", "1000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "max_udp_package_len", "1200"),
				),
			},
			{
				Config: testAccDayuDdosPolicyCaseUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyCaseExists(testDayuDdosPolicyCaseResourceKey),
					resource.TestCheckResourceAttrSet(testDayuDdosPolicyCaseResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "resource_type", "bgp"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "name", "tf_test_policy_case"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "platform_types.#", "2"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "web_api_urls.#", "2"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "app_protocols.#", "1"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_abroad", "no"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_initiate_tcp", "no"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_initiate_udp", "no"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "has_vpn", "no"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "tcp_start_port", "3000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "tcp_end_port", "4000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "udp_start_port", "1000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "udp_end_port", "2000"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "peer_tcp_port", "333"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "peer_udp_port", "111"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "tcp_footprint", "411"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "udp_footprint", "400"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "min_tcp_package_len", "900"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "max_tcp_package_len", "1100"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "min_udp_package_len", "900"),
					resource.TestCheckResourceAttr(testDayuDdosPolicyCaseResourceKey, "max_udp_package_len", "1100"),
				),
			},
		},
	})
}

func testAccCheckDayuDdosPolicyCaseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuDdosPolicyCaseResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of DDos policy case")
		}
		resourceType := items[0]
		sceneId := items[1]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		if err != nil {
			_, has, err = service.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete DDoS policy case %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuDdosPolicyCaseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of DDos policy case")
		}
		resourceType := items[0]
		sceneId := items[1]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		if err != nil {
			_, has, err = service.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("DDoS policy case %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuDdosPolicyCase string = `
resource "tencentcloud_dayu_ddos_policy_case" "test_policy_case" {
  resource_type         = "bgp"
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
`
const testAccDayuDdosPolicyCaseUpdate string = `
resource "tencentcloud_dayu_ddos_policy_case" "test_policy_case" {
  resource_type         = "bgp"
  name                  = "tf_test_policy_case"
  platform_types        = ["MOBILE","SERVER"]
  app_type              = "GAME"
  app_protocols         = ["all"]
  tcp_start_port		= "3000"
  tcp_end_port          = "4000"
  udp_start_port		= "1000"
  udp_end_port			= "2000"
  has_abroad			= "no"
  has_initiate_tcp		= "no"
  has_initiate_udp		= "no"
  peer_tcp_port			= "333"
  peer_udp_port			= "111"
  tcp_footprint		= "411"
  udp_footprint		= "400"
  web_api_urls			= ["abc.com/aaa.xls", "test.cn/bbb.png"]
  min_tcp_package_len	= "900"
  max_tcp_package_len	= "1100"
  min_udp_package_len	= "900"
  max_udp_package_len	= "1100"
  has_vpn				= "no"
}
`
