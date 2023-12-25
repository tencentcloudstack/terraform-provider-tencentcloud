package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesAuthAttachResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAuthAttach(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_create_discovery_anonymous_auth", "true"),
				),
			},
		},
	})
}

func testAccTkeAuthAttach() string {
	return tcacctest.TkeCIDRs + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.tke_cidr_a.1
  cluster_max_pod_num     = 32
  cluster_name            = "for-auth-attachment"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32
  cluster_os			  = "tlinux2.2(tkernel3)x86_64"

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                           = tencentcloud_kubernetes_cluster.managed_cluster.id
  issuer                               = ""
  auto_create_discovery_anonymous_auth = true
  use_tke_default                      = true
}
`
}
