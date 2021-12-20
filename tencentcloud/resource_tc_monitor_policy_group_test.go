package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudMonitorPolicyGroupResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorPolicyGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorPolicyGroupExists("tencentcloud_monitor_policy_group.group"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "group_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "policy_view_name", "cvm_device"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "remark", "this is a test policy group"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "is_union_rule", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.metric_id", "33"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.alarm_notify_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.alarm_notify_period", "600"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.calc_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.calc_value", "3"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.calc_period", "300"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.continue_period", "2"),
				),
			},
			{
				Config: testAccMonitorPolicyGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorPolicyGroupExists("tencentcloud_monitor_policy_group.group"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "group_name", "terraform_update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "policy_view_name", "cvm_device"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "remark", "this is a test policy group"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "is_union_rule", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.metric_id", "33"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.alarm_notify_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.alarm_notify_period", "600"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.calc_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.calc_value", "3"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.calc_period", "300"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.0.continue_period", "2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.metric_id", "30"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.alarm_notify_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.alarm_notify_period", "600"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.calc_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.calc_value", "30"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.calc_period", "300"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "conditions.1.continue_period", "2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.0.event_id", "39"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.0.alarm_notify_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.0.alarm_notify_period", "300"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.1.event_id", "40"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.1.alarm_notify_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_policy_group.group", "event_conditions.1.alarm_notify_period", "300"),
				),
			},
		},
	})
}

func testAccCheckMonitorPolicyGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_policy_group" {
			continue
		}
		groupIdStr := rs.Primary.ID

		if groupIdStr == "" {
			return fmt.Errorf("miss group_id[%v] ", groupIdStr)
		}
		groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
		if err != nil {
			return fmt.Errorf("id [%d] is broken", groupId)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		info, err := service.DescribePolicyGroup(ctx, groupId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				info, err = service.DescribePolicyGroup(ctx, groupId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}

		if info != nil {
			return fmt.Errorf("group %d found in DescribePolicyGroup", groupId)
		}
		log.Printf("[DEBUG]group %d delete ok", groupId)

	}
	return nil
}

func testAccCheckMonitorPolicyGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}
		groupIdStr := rs.Primary.ID

		if groupIdStr == "" {
			return fmt.Errorf("miss group_id[%v] ", groupIdStr)
		}
		groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
		if err != nil {
			return fmt.Errorf("id [%d] is broken", groupId)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		info, err := service.DescribePolicyGroup(ctx, groupId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				info, err = service.DescribePolicyGroup(ctx, groupId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if info == nil {
			return fmt.Errorf("group %d not found in DescribePolicyGroup", groupId)
		}

		return nil
	}
}

const testAccMonitorPolicyGroupBasic string = `
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "terraform_test"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 1
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
}
`
const testAccMonitorPolicyGroupUpdate string = `
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "terraform_update"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 0
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
  conditions {
    metric_id           = 30
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 2
    calc_value          = 30
    calc_period         = 300
    continue_period     = 2
  }
  event_conditions {
    event_id            = 39
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
  event_conditions {
    event_id            = 40
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
}
`
