package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTkeClusterEndpoint(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeClusterEndpointBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTkeExists("tencentcloud_kubernetes_cluster.managed_cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_internet", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_intranet", "true"),
				),
			},
			{
				ResourceName:            "tencentcloud_kubernetes_cluster_endpoint.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_intranet_subnet_id"},
			},
			{
				Config: testAccTkeClusterEndpointBasicUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTkeExists("tencentcloud_kubernetes_cluster.managed_cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_internet", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_endpoint.foo", "cluster_intranet", "true"),
				),
			},
		},
	})
}

const testAccTkeClusterEndpointBasicDeps = TkeCIDRs + TkeDataSource + ClusterAttachmentInstanceType + defaultImages + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_instances" "vpcs" {
  name = "keep_tke_exclusive_vpc"
}

data "tencentcloud_vpc_subnets" "sub" {
  vpc_id        = data.tencentcloud_vpc_instances.vpcs.instance_list.0.vpc_id
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "tf-auto-test-1-2"
  availability_zone = data.tencentcloud_vpc_subnets.sub.instance_list.0.availability_zone
  image_id          = var.default_img_id
  instance_type     = local.type1
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  vpc_id            = data.tencentcloud_vpc_instances.vpcs.instance_list.0.vpc_id
  subnet_id         =  data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  tags = data.tencentcloud_kubernetes_clusters.tke.list.0.tags # new added node will passive add tag by cluster
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

resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
  cluster_id  = tencentcloud_kubernetes_cluster.managed_cluster.id
  instance_id = tencentcloud_instance.foo.id
  password    = "Lo4wbdit"
  unschedulable = 0
}

resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  cluster_internet = true
  cluster_intranet = true
  managed_cluster_internet_security_policies = [
    "192.168.0.0/24"
  ]
  cluster_intranet_subnet_id = data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  depends_on = [
	tencentcloud_kubernetes_cluster_attachment.test_attach
  ]
}
`

const testAccTkeClusterEndpointBasic = testAccTkeClusterEndpointBasicDeps + `
resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  cluster_internet = true
  cluster_intranet = true
  managed_cluster_internet_security_policies = [
    "192.168.0.0/24"
  ]
  cluster_intranet_subnet_id = data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  depends_on = [
	tencentcloud_kubernetes_cluster_attachment.test_attach
  ]
}
`

const testAccTkeClusterEndpointBasicUpdate = testAccTkeClusterEndpointBasicDeps + `
resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  cluster_internet = false
  cluster_intranet = true
  cluster_intranet_subnet_id = data.tencentcloud_vpc_subnets.sub.instance_list.0.subnet_id
  depends_on = [
	tencentcloud_kubernetes_cluster_attachment.test_attach
  ]
}
`
