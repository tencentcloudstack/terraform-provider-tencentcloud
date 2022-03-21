package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudEksClusterCredentialDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudEksClusterCredentialBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_cluster_credential.cred", "addresses.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_cluster_credential.cred", "credential.ca_cert"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_cluster_credential.cred", "credential.token"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_cluster_credential.cred", "public_lb.0.enabled", "true"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_cluster_credential.cred", "public_lb.0.allow_from_cidrs.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_cluster_credential.cred", "public_lb.0.security_policies.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_cluster_credential.cred", "internal_lb.0.enabled", "true"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_cluster_credential.cred", "proxy_lb"),
				),
			},
		},
	})
}

const testAccTencentCloudEksClusterForCred = defaultVpcVariable + `
resource "tencentcloud_eks_cluster" "foo" {
  cluster_name = "tf-eks-test"
  k8s_version = "1.18.4"
  vpc_id = var.vpc_id
  subnet_ids = [
    var.subnet_id,
  ]
  cluster_desc = "test eks cluster created by terraform"
  service_subnet_id = var.subnet_id
  enable_vpc_core_dns = true
  internal_lb {
    enabled = true
    subnet_id = var.subnet_id
  }
  public_lb {
    enabled = true
    security_policies = ["192.168.1.1"]
  }
  tags = {
    test = "tf"
  }
}`

const testAccTencentCloudEksClusterCredentialBasic = testAccTencentCloudEksClusterForCred + `
data "tencentcloud_eks_cluster_credential" "cred" {
  cluster_id = tencentcloud_eks_cluster.foo.id
}
`
