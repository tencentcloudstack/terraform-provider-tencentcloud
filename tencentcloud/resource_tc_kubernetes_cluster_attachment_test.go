package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudTkeAttachResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAttachCluster(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeAttachExists("tencentcloud_kubernetes_cluster_attachment.test_attach"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_attachment.test_attach", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_attachment.test_attach", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckTkeAttachDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_kubernetes_cluster_attachment" {
			continue
		}
		clusterId := ""
		if items := strings.Split(rs.Primary.ID, "_"); len(items) != 2 {
			return fmt.Errorf("the resource id is corrupted")
		} else {
			clusterId = items[1]
		}

		_, has, err := service.DescribeCluster(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = service.DescribeCluster(ctx, clusterId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}

		if has {
			return fmt.Errorf("tke cluster attach delete fail,%s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckTkeAttachExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("tke cluster attach %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("tke cluster  attach id is not set")
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		instanceId, clusterId := "", ""

		if items := strings.Split(rs.Primary.ID, "_"); len(items) != 2 {
			return fmt.Errorf("the resource id is corrupted")
		} else {
			instanceId, clusterId = items[0], items[1]
		}

		service := TkeService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
				if err != nil {
					return retryError(err, InternalError)
				}
				return nil
			})
		}

		has := false
		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				has = true
			}
		}

		if !has {
			return fmt.Errorf("tke cluster attach cvm fail")
		}
		return nil

	}
}

func testAccTkeAttachCluster() string {

	return `

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

variable "default_instance_type" {
  default = "SA1.LARGE8"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}


data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["SA2"]
  }

  cpu_core_count = 8
  memory_size    = 16
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "tf-auto-test-1-1"
  availability_zone = var.availability_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = var.default_instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
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

resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
  cluster_id  = tencentcloud_kubernetes_cluster.managed_cluster.id
  instance_id = tencentcloud_instance.foo.id
  password    = "Lo4wbdit"
}
`
}
