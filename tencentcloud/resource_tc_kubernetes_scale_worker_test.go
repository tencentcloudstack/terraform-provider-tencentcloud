package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testTkeScaleWorkerResourceName = "tencentcloud_kubernetes_scale_worker"
var testTkeScaleWorkerResourceKey = testTkeScaleWorkerResourceName + ".test_scale"

func TestAccTencentCloudTkeScaleWorkerResource(t *testing.T) {
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
		ctx := context.WithValue(context.TODO(), "logId", logId)

		service := TkeService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() == "InternalError.ClusterNotFound" {
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
		ctx := context.WithValue(context.TODO(), "logId", logId)

		service := TkeService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() == "InternalError.ClusterNotFound" {
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

const testAccTkeScaleWorkerInstance string = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "vpc" {
  default = "` + defaultVpcId + `"
}

variable "subnet" {
  default = "` + defaultSubnetId + `"
}

variable "default_instance_type" {
  default = "SA1.LARGE8"
}

variable "scale_instance_type" {
  default = "S2.LARGE16"
}
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = "${var.vpc}"
  cluster_cidr            = "192.168.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.default_instance_type}"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = "${var.subnet}"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = "${tencentcloud_kubernetes_cluster.managed_cluster.id}"

  worker_config {
    count                      = 1
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.scale_instance_type}"
    subnet_id                  = "${var.subnet}"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 50
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "AABBccdd1122"
  }
}
`
