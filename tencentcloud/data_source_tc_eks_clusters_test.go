package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudEksClustersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudEKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEksClusterDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEksClusterExists("tencentcloud_eks_clusters.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_clusters.foo", "cluster_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "cluster_name", "tf-eks-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_clusters.foo", "list.0.cluster_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.cluster_name", "tf-eks-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.k8s_version", "1.18.4"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_clusters.foo", "list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.dns_servers.0.domain", "example2.org"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.dns_servers.0.servers.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.dns_servers.0.servers.0", "10.0.0.1:80"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.dns_servers.0.servers.1", "10.0.0.1:81"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.cluster_desc", "test eks cluster created by terraform"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eks_clusters.foo", "list.0.service_subnet_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.enable_vpc_core_dns", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.need_delete_cbs", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_eks_clusters.foo", "list.0.tags.0.test", "tf"),
				),
			},
		},
	})
}

func testAccCheckEksClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		var (
			has bool
			err error
		)
		if !ok {
			return fmt.Errorf("EKS cluster %s is not found", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("EKS cluster id is not set")
		}
		eksService := EksService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, has, err = eksService.DescribeEksCluster(ctx, id)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = eksService.DescribeEksCluster(ctx, id)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("EKS cluster is not found")
		}
		return nil
	}
}

const testAccEksClusterDataSource = defaultVpcVariable + `
resource "tencentcloud_eks_cluster" "foo" {
 cluster_name = "tf-eks-test"
 k8s_version = "1.18.4"
 vpc_id = var.vpc_id
 subnet_ids = [
   var.subnet_id,
 ]
 dns_servers {
   domain = "example2.org"
   servers = ["10.0.0.1:80", "10.0.0.1:81"]
 }
 cluster_desc = "test eks cluster created by terraform"
 service_subnet_id = var.subnet_id
 enable_vpc_core_dns = true
 need_delete_cbs = true
 tags = {
   test = "tf"
 }
}

data "tencentcloud_eks_clusters" "foo" {
  cluster_id = tencentcloud_eks_cluster.foo.id
}
`
