package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

var testInternationalTkeClusterName = "tencentcloud_kubernetes_cluster"
var testInternationalTkeClusterResourceKey = testInternationalTkeClusterName + ".managed_cluster"

func TestAccTencentCloudInternationalKubernetesResource_cluster(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTkeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalTkeCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_cidr", "10.31.0.0/23"),
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
					//resource.TestCheckResourceAttr(testTkeClusterResourceKey, "security_policy.#", "2"),
					//resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "cluster_external_endpoint"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_level", "L5"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "auto_upgrade_cluster_level", "true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "labels.test1", "test1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "labels.test2", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_internet_domain", "tf.cluster-internet.com"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_intranet_domain", "tf.cluster-intranet.com"),
				),
			},
		},
	})
}

const InternationalTkeDeps = tcacctest.TkeExclusiveNetwork + tcacctest.TkeInstanceType + tcacctest.TkeCIDRs + tcacctest.DefaultImages + tcacctest.DefaultSecurityGroupData

const testAccInternationalTkeCluster = InternationalTkeDeps + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_a.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  cluster_internet_domain                    = "tf.cluster-internet.com"
  cluster_intranet                           = true
  cluster_intranet_domain                    = "tf.cluster-intranet.com"
  cluster_version                            = "1.22.5"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L5"
  auto_upgrade_cluster_level				 = true
  cluster_intranet_subnet_id                 = local.subnet_id
  cluster_internet_security_group               = local.sg_id
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]

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
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "test" = "test"
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
`
