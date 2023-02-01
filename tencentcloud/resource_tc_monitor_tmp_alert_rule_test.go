package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorAlertRuleResource_basic -v
func TestAccTencentCloudMonitorAlertRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlertRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertRuleExists("tencentcloud_monitor_tmp_alert_rule.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "rule_name", "test-rule_name"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "receivers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "expr", "increase(mysql_global_status_slow_queries[1m]) > 0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "duration", "4m"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "labels.0.key", "hello1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "labels.0.value", "world1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "annotations.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "annotations.0.key", "hello2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "annotations.0.value", "world2"),
				),
			},
			{
				Config: testAlertRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertRuleExists("tencentcloud_monitor_tmp_alert_rule.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "rule_name", "test-rule_name_update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "receivers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "expr", "increase(mysql_global_status_slow_queries[1m]) > 1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "duration", "2m"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "labels.0.key", "hello3"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "labels.0.value", "world3"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "annotations.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "annotations.0.key", "hello4"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_alert_rule.basic", "annotations.0.value", "world4"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_alert_rule.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAlertRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_alert_rule" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instance, err := service.DescribeMonitorTmpAlertRuleById(ctx, ids[0], ids[1])
		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAlertRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		instance, err := service.DescribeMonitorTmpAlertRuleById(ctx, ids[0], ids[1])
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAlertRuleVar = `
variable "prometheus_id" {
  default = "` + defaultPrometheusId + `"
}
`
const testAlertRule_basic = testAlertRuleVar + `
resource "tencentcloud_monitor_tmp_alert_rule" "basic" {
  instance_id	= var.prometheus_id
  rule_name		= "test-rule_name"
  receivers 	= ["notice-tj75hgqj"]
  expr			= "increase(mysql_global_status_slow_queries[1m]) > 0"
  duration	    = "4m"
  rule_state	= 2
  labels {
    key   = "hello1"
    value = "world1"
  }
  annotations {
    key   = "hello2"
    value = "world2"
  }
}`

const testAlertRule_update = testAlertRuleVar + `
resource "tencentcloud_monitor_tmp_alert_rule" "basic" {
  instance_id	= var.prometheus_id
  rule_name		= "test-rule_name_update"
  receivers 	= ["notice-tj75hgqj"]
  expr			= "increase(mysql_global_status_slow_queries[1m]) > 1"
  duration	    = "2m"
  rule_state	= 2
  labels {
    key   = "hello3"
    value = "world3"
  }
  annotations {
    key   = "hello4"
    value = "world4"
  }
}`
