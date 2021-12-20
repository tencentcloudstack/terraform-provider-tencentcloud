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

func TestAccTencentCloudMonitorBindingReceiverResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorBindingReceiverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorBindingReceiverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorBindingReceiverExists("tencentcloud_monitor_binding_receiver.receiver"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_binding_receiver.receiver", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.start_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.end_time", "86399"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.receiver_type", "group"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.receive_language", "en-US"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.notify_way.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.notify_way.0", "SMS"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_binding_receiver.receiver", "receivers.0.receiver_group_list.#", "1"),
				),
			},
		},
	})
}

func testAccCheckMonitorBindingReceiverDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_binding_receiver" {
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

func testAccCheckMonitorBindingReceiverExists(n string) resource.TestCheckFunc {
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
		if len(info.ReceiverInfos) > 0 {
			return nil
		} else {
			return fmt.Errorf("group %d not found receiver in DescribePolicyGroup", groupId)
		}
	}
}

const testAccMonitorBindingReceiverInstance string = `

data "tencentcloud_cam_groups" "groups" {

}

resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "terraform_test"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
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
resource "tencentcloud_monitor_binding_receiver" "receiver" {
  group_id = tencentcloud_monitor_policy_group.group.id
  receivers {
    start_time          = 0
    end_time            = 86399
    notify_way          = ["SMS"]
    receiver_type       = "group"
    receiver_group_list = [data.tencentcloud_cam_groups.groups.group_list[0].group_id]
    receive_language    = "en-US"
  }
}
`
