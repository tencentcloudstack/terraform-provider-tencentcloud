package dayu_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcdayu "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"
)

var testDayuL7RuleResourceName = "tencentcloud_dayu_l7_rule"
var testDayuL7RuleResourceKey = testDayuL7RuleResourceName + ".test_rule"

func TestAccTencentCloudDayuL7RuleResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDayuL7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuL7Rule, tcacctest.DefaultDayuBgpIp, tcacctest.DefaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL7RuleExists(testDayuL7RuleResourceKey),
					resource.TestCheckResourceAttrSet(testDayuL7RuleResourceKey, "rule_id"),
					resource.TestCheckResourceAttrSet(testDayuL7RuleResourceKey, "status"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "name", "rule_test"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "source_type", "2"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "source_list.#", "2"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "switch", "true"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "protocol", "https"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "ssl_id", tcacctest.DefaultSshCertificate),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_code", "31"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_switch", "true"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_interval", "30"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_path", "/"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_method", "GET"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_health_num", "5"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "health_check_unhealth_num", "10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDayuL7RuleUpdate, tcacctest.DefaultDayuBgpIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuL7RuleExists(testDayuL7RuleResourceKey),
					testAccCheckDayuL7RuleExists(testDayuL7RuleResourceKey),
					resource.TestCheckResourceAttrSet(testDayuL7RuleResourceKey, "rule_id"),
					resource.TestCheckResourceAttrSet(testDayuL7RuleResourceKey, "status"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "name", "rule_test"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "source_type", "1"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "source_list.#", "1"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "switch", "false"),
					resource.TestCheckResourceAttr(testDayuL7RuleResourceKey, "protocol", "http"),
				),
			},
		},
	})
}

func testAccCheckDayuL7RuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuL7RuleResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of L7 rule")
		}
		resourceType := items[0]
		resourceId := items[1]
		ruleId := items[2]

		service := svcdayu.NewDayuService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, _, has, err := service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			_, _, has, err = service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete L7 rule %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuL7RuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of L7 rule")
		}
		resourceType := items[0]
		resourceId := items[1]
		ruleId := items[2]

		service := svcdayu.NewDayuService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, _, has, err := service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			_, _, has, err = service.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("L7 rule %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuL7Rule string = `
resource "tencentcloud_dayu_l7_rule" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  name					= "rule_test"
  domain				= "zhaoshaona.com"
  protocol				= "https"
  switch				= true
  source_type			= 2
  source_list 			= ["1.1.1.1:80","2.2.2.2"]
  ssl_id				= "%s"
  health_check_switch	= true
  health_check_code		= 31
  health_check_interval = 30
  health_check_method	= "GET"
  health_check_path		= "/"
  health_check_health_num = 5
  health_check_unhealth_num = 10
}
`
const testAccDayuL7RuleUpdate string = `
resource "tencentcloud_dayu_l7_rule" "test_rule" {
  resource_type         = "bgpip"
  resource_id 			= "%s"
  name					= "rule_test"
  domain				= "zhaoshaona.com"
  protocol				= "http"
  switch				= false
  source_type			= 1
  source_list 			= ["zhaoshaona.com"]
}
`
