package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTkeClusterAsName = "tencentcloud_kubernetes_as_scaling_group"
var testTkeClusterAsResourceKey = testTkeClusterAsName + ".as_test"

func TestAccTencentCloudTkeAsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeAsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAsCluster(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeAsExists,
					resource.TestCheckResourceAttrSet(testTkeClusterAsResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeClusterAsResourceKey, "auto_scaling_group.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterAsResourceKey, "auto_scaling_config.#", "1"),
				),
			},
		},
	})
}

func testAccCheckTkeAsDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	rs, ok := s.RootModule().Resources[testTkeClusterAsResourceKey]
	if !ok {
		return fmt.Errorf("tke as group %s is not found", testTkeClusterAsResourceKey)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("tke  as group  id is not set")
	}
	items := strings.Split(rs.Primary.ID, ":")
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_as_scaling_group id %s is broken", rs.Primary.ID)
	}
	asGroupId := items[1]

	_, has, err := service.DescribeAutoScalingGroupById(ctx, asGroupId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	} else {
		return fmt.Errorf("tke as group %s still exist", asGroupId)
	}

}

func testAccCheckTkeAsExists(s *terraform.State) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	rs, ok := s.RootModule().Resources[testTkeClusterAsResourceKey]
	if !ok {
		return fmt.Errorf("tke as group %s is not found", testTkeClusterAsResourceKey)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("tke  as group  id is not set")
	}

	items := strings.Split(rs.Primary.ID, ":")
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_as_scaling_group id  %s is broken", rs.Primary.ID)
	}
	asGroupId := items[1]

	_, has, err := service.DescribeAutoScalingGroupById(ctx, asGroupId)
	if err != nil {
		return err
	}
	if has == 1 {
		return nil
	} else {
		return fmt.Errorf("tke as group %s query fail.", asGroupId)
	}

}

func testAccTkeAsCluster() string {
	return fmt.Sprintf(`

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
  default = "SN3ne.8XLARGE64"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf-tke-unit-test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

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
}


resource "tencentcloud_kubernetes_as_scaling_group" "as_test" {

  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id

  auto_scaling_group {
    scaling_group_name   = "tf-tke-as-group-unit-test"
    max_size             = "5"
    min_size             = "0"
    vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
    subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
    project_id           = 0
    default_cooldown     = 400
    desired_capacity     = "0"
    termination_policies = ["NEWEST_INSTANCE"]
    retry_policy         = "INCREMENTAL_INTERVALS"

    tags = {
      "test" = "test"
    }

  }

  auto_scaling_config {
    configuration_name = "tf-tke-as-config-unit-test"
    instance_type      = var.default_instance_type
    project_id         = 0
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"

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

    instance_tags = {
      tag = "as"
    }

  }
}

`,
	)
}
