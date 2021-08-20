package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTkeClusterName = "tencentcloud_kubernetes_cluster"
var testTkeClusterResourceKey = testTkeClusterName + ".managed_cluster"

func TestAccTencentCloudTkeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeCluster("test", "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_cidr", "10.31.0.0/16"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_max_pod_num", "32"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_name", "test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_desc", "test cluster desc"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_node_num", "1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "worker_instances_list.#", "1"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "worker_instances_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "certification_authority"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "user_name"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "password"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "tags.test", "test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "security_policy.#", "2"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "cluster_external_endpoint"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "labels.test1", "test1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "labels.test2", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "unschedulable", "0"),
				),
			},
			{
				Config: testAccTkeCluster("abc", "abc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckNoResourceAttr(testTkeClusterResourceKey, "tags.test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "tags.abc", "abc"),
				),
			},
		},
	})
}

func testAccCheckTkeDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTkeClusterName {
			continue
		}
		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}

		if !has {
			log.Printf("[DEBUG]tke cluster  %s delete  ok", rs.Primary.ID)
			return nil
		} else {
			return fmt.Errorf("tke cluster delete fail,%s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckTkeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("tke cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("tke cluster id is not set")
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := TkeService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}
		if !has {
			return fmt.Errorf("tke cluster create fail")
		} else {
			log.Printf("[DEBUG]tke cluster  %s create  ok", rs.Primary.ID)
			return nil
		}

	}
}

func testAccTkeCluster(key, value string) string {
	return fmt.Sprintf(`
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "10.31.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr                               = var.cluster_cidr
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  cluster_version                            = "1.18.4"
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
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
    img_id                     = "img-rkiynh11"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      file_system = "ext3"
	  auto_format_and_mount = "true"
	  mount_target = "/var/lib/docker"
      disk_partition = "/dev/sdb1"
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
	cam_role_name			= "CVM_QcsRole"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "%s" = "%s"
  }

  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
}
`, key, value,
	)
}
