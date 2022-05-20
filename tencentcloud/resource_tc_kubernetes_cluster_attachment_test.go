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
	t.Parallel()
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
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_attachment.test_attach", "unschedulable"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_attachment.test_attach", "labels.test1", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_attachment.test_attach", "labels.test2", "test2"),
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
		instanceId := ""
		clusterId := ""
		if items := strings.Split(rs.Primary.ID, "_"); len(items) != 2 {
			return fmt.Errorf("the resource id is corrupted")
		} else {
			instanceId = items[0]
			clusterId = items[1]
		}

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}

		var instanceStillExists bool

		for i := range workers {
			id := workers[i].InstanceId
			if instanceId == id {
				instanceStillExists = true
				break
			}
		}

		if instanceStillExists {
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

const ClusterAttachmentInstanceType = `
data "tencentcloud_instance_types" "ins_type" {
  availability_zone = "ap-guangzhou-3"
  cpu_core_count    = 2
  memory_size       = 2
  exclude_sold_out  = true
}

locals {
  type1 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"][0].instance_type
  type2 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"][1].instance_type
}
`

func testAccTkeAttachCluster() string {

	return TkeDataSource + ClusterAttachmentInstanceType + defaultImages + `
variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

data "tencentcloud_vpc_instances" "vpcs" {
  name = "keep_tke_exclusive_vpc"
}

data "tencentcloud_vpc_subnets" "sub" {
  vpc_id        = data.tencentcloud_vpc_instances.vpcs.instance_list.0.vpc_id
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "tf-auto-test-1-1"
  availability_zone = data.tencentcloud_vpc_subnets.sub.instance_list.0.availability_zone
  image_id          = var.default_img_id
  instance_type     = local.type1
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  vpc_id            = data.tencentcloud_vpc_instances.vpcs.instance_list.0.vpc_id
  subnet_id         =  data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  tags = data.tencentcloud_kubernetes_clusters.tke.list.0.tags # new added node will passive add tag by cluster
}

resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
  cluster_id  = local.cluster_id
  instance_id = tencentcloud_instance.foo.id
  password    = "Lo4wbdit"
  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

}
`
}
