package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoApplicationProxyRule_basic -v
func TestAccTencentCloudTeoApplicationProxyRule_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_application_proxy_rule" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]
		ruleId := idSplit[2]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTeoApplicationProxyRule = testAccTeoApplicationProxy + `

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
  proxy_id          = tencentcloud_teo_application_proxy.basic.proxy_id
  session_persist   = false
  status            = "online"
  zone_id           = tencentcloud_teo_zone.basic.id
}

`
