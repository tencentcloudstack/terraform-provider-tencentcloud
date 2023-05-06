package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoApplicationProxyRule_basic -v
func TestAccTencentCloudTeoApplicationProxyRule_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationProxyRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxyRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationProxyRuleExists("tencentcloud_teo_application_proxy_rule.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "forward_client_ip", "TOA"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "origin_type", "custom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "origin_port", "8083"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "origin_value.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "port.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "proto", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "session_persist", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy_rule.basic", "status", "online"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy_rule.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckApplicationProxyRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_application_proxy_rule" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]
		ruleId := idSplit[2]

		agents, err := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if agents != nil {
			return fmt.Errorf("zone ApplicationProxyRule %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckApplicationProxyRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]
		ruleId := idSplit[2]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoApplicationProxyRule(ctx, zoneId, proxyId, ruleId)
		if agents == nil {
			return fmt.Errorf("zone ApplicationProxyRule %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoApplicationProxyRuleVar = `
variable "default_zone_id" {
  default = "` + defaultZoneId + `"
}
variable "proxy_id" {
  default = "` + applicationProxyId + `"
}
`

const testAccTeoApplicationProxyRule = testAccTeoApplicationProxyRuleVar + `

resource "tencentcloud_teo_application_proxy_rule" "basic" {
  forward_client_ip = "TOA"
  origin_type       = "custom"
  origin_port       = "8083"
  origin_value      = [
    "127.0.0.1",
  ]
  port              = [
    "8083",
  ]
  proto             = "TCP"
  proxy_id          = var.proxy_id
  session_persist   = false
  status            = "online"
  zone_id           = var.default_zone_id
}

`
