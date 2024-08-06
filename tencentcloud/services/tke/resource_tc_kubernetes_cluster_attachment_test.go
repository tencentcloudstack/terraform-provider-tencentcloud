package tke_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctke "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudKubernetesClusterAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctke.NewTkeService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
				if err != nil {
					return tccommon.RetryError(err)
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

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		instanceId, clusterId := "", ""

		if items := strings.Split(rs.Primary.ID, "_"); len(items) != 2 {
			return fmt.Errorf("the resource id is corrupted")
		} else {
			instanceId, clusterId = items[0], items[1]
		}

		service := svctke.NewTkeService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
				if err != nil {
					return tccommon.RetryError(err, tccommon.InternalError)
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

	return TkeNewDeps + `
variable "cluster_cidr" {
  default = "10.31.0.0/16"
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_instances" "vpcs" {
  name = "keep_tke_exclusive_vpc"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = local.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_internet        = false # (can be ignored) open it after the nodes added
  cluster_version         = "1.22.5"
  cluster_os              = "tlinux2.2(tkernel3)x86_64"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  # without any worker config
}

resource "tencentcloud_instance" "foo_attachment" {
  instance_name     = "tf-auto-test-1-1"
  availability_zone = var.availability_zone
  image_id          = var.default_img_id
  instance_type     = local.final_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  vpc_id            = local.vpc_id
  subnet_id         = local.subnet_id1
}

resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
  cluster_id  = tencentcloud_kubernetes_cluster.example.id
  instance_id = tencentcloud_instance.foo_attachment.id
  password    = "Lo4wbdit"
  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

}
`
}
