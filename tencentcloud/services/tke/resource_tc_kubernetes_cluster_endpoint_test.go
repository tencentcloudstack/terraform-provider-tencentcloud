package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesClusterEndpointResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeClusterEndpointBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTkeExists("tencentcloud_kubernetes_cluster.managed_cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_internet", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_intranet", "true"),
					//resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "managed_cluster_internet_security_policies.#", "1"),
					//resource.TestCheckResourceAttr(
					//	"tencentcloud_kubernetes_cluster_endpoint.foo",
					//	"managed_cluster_internet_security_policies.0",
					//	"192.168.0.0/24",
					//),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_intranet_subnet_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_cluster_endpoint.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// FIXME waiting for DescribeEndpoints available
					"cluster_intranet_subnet_id",
					"cluster_internet",
					"cluster_internet_security_group",
					"cluster_intranet",
					"managed_cluster_internet_security_policies",
				},
			},
			{
				Config: testAccTkeClusterEndpointBasicUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTkeExists("tencentcloud_kubernetes_cluster.managed_cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_internet", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "extensive_parameters"),
				),
			},
			{
				Config: testAccTkeClusterEndpointBasicUpdate2,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTkeExists("tencentcloud_kubernetes_cluster.managed_cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_internet", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_intranet", "true"),
				),
			},
		},
	})
}

// const testAccTkeClusterEndpointNewSG = `
// data "tencentcloud_security_groups" "new_sg" {
//   name = "keep-tke-ep-sg-fwf8zdkx"
// }

// locals {
//   new_sg = data.tencentcloud_security_groups.new_sg.security_groups.0.security_group_id
// }

// `

const testAccTkeClusterEndpointBasicDeps = tcacctest.TkeCIDRs +
	tcacctest.TkeDataSource +
	tcacctest.TkeDefaultNodeInstanceVar +
	tcacctest.DefaultImages +
	tcacctest.DefaultSecurityGroupData +
	`
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_instances" "vpcs" {
  name = "keep_tke_exclusive_vpc"
}

data "tencentcloud_vpc_subnets" "sub" {
  vpc_id        = data.tencentcloud_vpc_instances.vpcs.instance_list.0.vpc_id
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.sub.instance_list.0.vpc_id
  cluster_cidr            = var.tke_cidr_a.3
  cluster_max_pod_num     = 32
  cluster_name            = "for-endpoint"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32
  cluster_os			  = "tlinux2.2(tkernel3)x86_64"

  cluster_deploy_type = "MANAGED_CLUSTER"
}

data "tencentcloud_security_groups" "sg" {   
  name = "default"
}

locals {
  new_cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
}

resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "test-endpoint-attachment"
  cluster_id = local.new_cluster_id
  max_size = 1
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.sub.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 1
  enable_auto_scale    = true
  scaling_group_name	   = "basic_group"
  default_cooldown		   = 400
  termination_policies	   = ["OLDEST_INSTANCE"]

  auto_scaling_config {
    instance_type      = var.ins_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    orderly_security_group_ids = [local.sg_id]

    cam_role_name = "TCB_QcsRole"
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
  node_os="tlinux2.2(tkernel3)x86_64"
}

`

const testAccTkeClusterEndpointBasic = testAccTkeClusterEndpointBasicDeps + `
resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = local.new_cluster_id
  cluster_internet = true
  cluster_intranet = true
  cluster_internet_security_group = local.sg_id
  managed_cluster_internet_security_policies = [
    "192.168.0.0/24"
  ]
  cluster_intranet_subnet_id = data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  depends_on = [
	tencentcloud_kubernetes_node_pool.np_test
  ]
}
`

const testAccTkeClusterEndpointBasicUpdate = testAccTkeClusterEndpointBasicDeps + `
resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = local.new_cluster_id
  cluster_internet = true
  cluster_intranet = true
  cluster_internet_security_group = local.sg_id2
  cluster_intranet_subnet_id = data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  extensive_parameters = jsonencode({
    InternetAccessible = {
      InternetChargeType = "BANDWIDTH_POSTPAID_BY_HOUR"
      InternetMaxBandwidthOut = 10
    }
  })
  depends_on = [
	tencentcloud_kubernetes_node_pool.np_test
  ]
}
`

const testAccTkeClusterEndpointBasicUpdate2 = testAccTkeClusterEndpointBasicDeps + `
resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = local.new_cluster_id
  cluster_internet = false
  cluster_intranet = true
  cluster_internet_security_group = local.sg_id2
  cluster_intranet_subnet_id = data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  depends_on = [
	tencentcloud_kubernetes_node_pool.np_test
  ]
}
`
