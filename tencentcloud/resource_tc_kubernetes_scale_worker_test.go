package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var testTkeScaleWorkerResourceName = "tencentcloud_kubernetes_scale_worker"
var testTkeScaleWorkerResourceKey = testTkeScaleWorkerResourceName + ".test_scale"

func TestAccTencentCloudTkeScaleWorkerResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeScaleWorkerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeScaleWorkerInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeScaleWorkerExists(testTkeScaleWorkerResourceKey),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "cluster_id"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_config.#"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "worker_instances_list.#", "1"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_instances_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_instances_list.0.instance_role"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "unschedulable"),
				),
			},
		},
	})
}

func testAccCheckTkeScaleWorkerDestroy(s *terraform.State) error {
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
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

func testAccCheckTkeScaleWorkerExists(n string) resource.TestCheckFunc {
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
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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

const testAccTkeScaleWorkerInstanceBasic = defaultAzVariable + TkeExclusiveNetwork + TkeDataSource

const testAccTkeScaleWorkerInstance string = testAccTkeScaleWorkerInstanceBasic + `
variable "scale_instance_type" {
  default = "S2.LARGE8"
}

resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = local.cluster_id
  
  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]	

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  unschedulable = 0

  worker_config {
    count                      				= 1
    availability_zone          				= var.default_az
    instance_type              				= var.scale_instance_type
    subnet_id                  				= local.subnet_id
    system_disk_type           				= "CLOUD_SSD"
    system_disk_size           				= 50
    internet_charge_type       				= "TRAFFIC_POSTPAID_BY_HOUR"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service 				= false
    enhanced_monitor_service  				= false
    user_data                 				= "dGVzdA=="
    password                  				= "AABBccdd1122"
  }
}
`
