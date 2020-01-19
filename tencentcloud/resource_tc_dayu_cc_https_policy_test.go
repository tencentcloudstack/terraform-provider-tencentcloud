package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testDayuCCHttpsPolicyResourceName = "tencentcloud_dayu_cc_https_policy"
var testDayuCCHttpsPolicyResourceKey = testDayuCCHttpsPolicyResourceName + ".test_policy"

func TestAccTencentCloudDayuCCHttpsPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpsPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuCCHttpsPolicy, defaultDayuBgpIp, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpsPolicyExists(testDayuCCHttpsPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "rule_id"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "domain"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "name", "policy_test"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "action", "drop"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "rule_list.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDayuCCHttpsPolicyUpdate, defaultDayuBgpIp, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpsPolicyExists(testDayuCCHttpsPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "rule_id"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpsPolicyResourceKey, "domain"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "name", "policy_test"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "action", "alg"),
					resource.TestCheckResourceAttr(testDayuCCHttpsPolicyResourceKey, "rule_list.#", "1"),
				),
			},
		},
	})
}

func testAccCheckDayuCCHttpsPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuCCHttpsPolicyResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of DDos policy case")
		}
		resourceType := items[0]
		resourceId := items[1]
		policyId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if err != nil {
			_, has, err = service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete CC https policy %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuCCHttpsPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of DDos policy case")
		}
		resourceType := items[0]
		resourceId := items[1]
		policyId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if err != nil {
			_, has, err = service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("CC https policy %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuCCHttpsPolicy string = `
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

resource "tencentcloud_dayu_cc_https_policy" "test_policy" {
  resource_type         = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id 			= tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id				= tencentcloud_dayu_l7_rule.test_rule.rule_id
  domain				= tencentcloud_dayu_l7_rule.test_rule.domain
  name					= "policy_test"
  action				= "drop"
  switch				= true

  rule_list {
	skey 				= "cgi"
	operator			= "include"
	value				= "123"
	}
}
`
const testAccDayuCCHttpsPolicyUpdate string = `
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

resource "tencentcloud_dayu_cc_https_policy" "test_policy" {
  resource_type         = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id 			= tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id				= tencentcloud_dayu_l7_rule.test_rule.rule_id
  domain				= tencentcloud_dayu_l7_rule.test_rule.domain
  name					= "policy_test"
  action				= "alg"
  switch				= true

  rule_list {
	skey 				= "ua"
	operator			= "equal"
	value				= "111"
	}
}
`
