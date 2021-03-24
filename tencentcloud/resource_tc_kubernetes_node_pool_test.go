package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTkeClusterNodePoolName = "tencentcloud_kubernetes_node_pool"
var testTkeClusterNodePoolResourceKey = testTkeClusterNodePoolName + ".np_test"

func TestAccTencentCloudTkeNodePoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeNodePoolCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeNodePoolExists,
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "taints.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "labels.test1", "test1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "labels.test2", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "max_size", "6"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "min_size", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "desired_capacity", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "name", "mynodepool"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "unschedulable", "0"),
				),
			},
			{
				Config: testAccTkeNodePoolClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeNodePoolExists,
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "max_size", "5"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "min_size", "2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "labels.test3", "test3"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "desired_capacity", "2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "name", "mynodepoolupdate"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_os", "ubuntu18.04.1x86_64"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "unschedulable", "1"),
				),
			},
		},
	})
}

func testAccCheckTkeNodePoolDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	rs, ok := s.RootModule().Resources[testTkeClusterNodePoolResourceKey]
	if !ok {
		return fmt.Errorf("tke node pool %s is not found", testTkeClusterNodePoolResourceKey)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("tke  node pool id is not set")
	}
	items := strings.Split(rs.Primary.ID, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id %s is broken", rs.Primary.ID)
	}
	clusterId := items[0]
	nodePoolId := items[1]

	_, has, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return fmt.Errorf("tke node pool %s still exist", nodePoolId)
	}

}

func testAccCheckTkeNodePoolExists(s *terraform.State) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	rs, ok := s.RootModule().Resources[testTkeClusterNodePoolResourceKey]
	if !ok {
		return fmt.Errorf("tke node pool %s is not found", testTkeClusterNodePoolResourceKey)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("tke node pool id is not set")
	}

	items := strings.Split(rs.Primary.ID, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  %s is broken", rs.Primary.ID)
	}
	clusterId := items[0]
	nodePoolId := items[1]

	_, has, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		return err
	}
	if has {
		return nil
	} else {
		return fmt.Errorf("tke node pool %s query fail.", nodePoolId)
	}

}

const testAccTkeNodePoolClusterBasic = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf-tke-unit-test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

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
}`

const testAccTkeNodePoolCluster string = testAccTkeNodePoolClusterBasic + `
resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "mynodepool"
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size = 6
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 1
  enable_auto_scale    = true

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = ["sg-24vswocp"]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

  }
  unschedulable = 0
  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
	key = "test_taint"
    value = "taint_value"
    effect = "PreferNoSchedule"
  }

  node_config {
      extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
  }
}

`

const testAccTkeNodePoolClusterUpdate string = testAccTkeNodePoolClusterBasic + `
resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "mynodepoolupdate"
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size = 5
  min_size = 2
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 2
  enable_auto_scale    = false
  node_os = "ubuntu18.04.1x86_64"
  delete_keep_instance = true

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = ["sg-24vswocp"]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

  }
  unschedulable = 1
  labels = {
    "test3" = "test3",
    "test2" = "test2",
  }

  node_config {
      extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
  }
}
`
