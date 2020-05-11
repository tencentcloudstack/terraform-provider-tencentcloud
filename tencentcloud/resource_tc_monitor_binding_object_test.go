package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudMonitorBindingObjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorBindingObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeScaleWorkerInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorBindingObjectExists(testTkeScaleWorkerResourceKey),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "cluster_id"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_config.#"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "worker_instances_list.#", "1"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_instances_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_instances_list.0.instance_role")),
			},
		},
	})
}

func testAccCheckMonitorBindingObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTkeScaleWorkerResourceName {
			continue
		}
		instanceId := rs.Primary.Attributes["worker_instances_list.0.instance_id"]
		clusterId := rs.Primary.Attributes["cluster_id"]

		if clusterId == "" || instanceId == "" {
			return fmt.Errorf("miss worker_instances_list.0.instance_id[%s] or cluster_id[%s]", instanceId, clusterId)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		service := TkeService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() == "InvalidParameter.ClusterNotFound" {
						return nil
					}
				}
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}

		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				return fmt.Errorf("cvm %s found in DescribeClusterInstances", instanceId)
			}
		}
		log.Printf("[DEBUG]instance %s delelte ok", instanceId)

	}
	return nil
}

func testAccCheckMonitorBindingObjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("tke worker scale instance %s is not found", n)
		}
		instanceId := rs.Primary.Attributes["worker_instances_list.0.instance_id"]
		clusterId := rs.Primary.Attributes["cluster_id"]

		if clusterId == "" || instanceId == "" {
			return fmt.Errorf("miss worker_instances_list.0.instance_id[%s] or cluster_id[%s]", instanceId, clusterId)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		service := TkeService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() == "InvalidParameter.ClusterNotFound" {
						return nil
					}
				}

				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}

		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				log.Printf("[DEBUG]instance %s create ok", instanceId)
				return nil
			}
		}
		return fmt.Errorf("cvm %s not found in DescribeClusterInstances", instanceId)
	}
}

const testAccMonitorBindingObjectInstance string = `


resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
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
resource "tencentcloud_monitor_binding_object" "binding" {
  group_id = tencentcloud_monitor_policy_group.group.id
  dimensions {
    dimensions_json = "{\"unInstanceId\":\"ins-lqm9lahu\"}"
  }
}
`
