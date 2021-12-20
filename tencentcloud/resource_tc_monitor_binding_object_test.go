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

func TestAccTencentCloudMonitorBindingObjectResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorBindingObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorBindingObjectInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorBindingObjectExists("tencentcloud_monitor_binding_object.binding"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_binding_object.binding", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_object.binding", "dimensions.#", "1"),
				),
			},
		},
	})
}

func testAccCheckMonitorBindingObjectDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_binding_object" {
			continue
		}
		groupIdStr := rs.Primary.Attributes["group_id"]

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
		log.Printf("[DEBUG]group and receivers %d delete ok", groupId)

	}
	return nil
}

func testAccCheckMonitorBindingObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}
		groupIdStr := rs.Primary.Attributes["group_id"]

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

		objects, err := service.DescribeBindingPolicyObjectList(ctx, groupId)

		if err != nil {
			return err
		}
		if len(objects) < 1 {
			return fmt.Errorf("group %d binding object fail", groupId)
		}
		return nil
	}
}

const testAccMonitorBindingObjectInstance string = `
data "tencentcloud_instances" "instances" {
}
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

resource "tencentcloud_monitor_binding_object" "binding" {
  group_id = tencentcloud_monitor_policy_group.group.id
  dimensions {
    dimensions_json = "{\"unInstanceId\":\"${data.tencentcloud_instances.instances.instance_list[0].instance_id}\"}"
  }
}
`
