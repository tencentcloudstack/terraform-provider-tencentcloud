package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccTencentCloudEKSCluster_basic(t *testing.T) {
	t.Parallel()

	fmt.Printf("file====:\n$%s", testAccEksCluster)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudEKSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEksCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccTencentCloudEKSClusterExists("resource_eks_cluster.foo"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "cluster_name", "tf-ci-test"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "k8s_version", "1.12.8"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "cluster_desc", "test eks cluster created by terraform"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.#", "1"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.domain", "example2.org"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.servers.#", "2"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.servers.0", "10.0.0.1:80"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.servers.1", "10.0.0.1:81"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "enable_vpc_core_dns", "true"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "need_delete_cbs", "true"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "tags.hello", "world"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet("resource_eks_cluster.foo", "subnet_ids.0"),
					resource.TestCheckResourceAttrSet("resource_eks_cluster.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("resource_eks_cluster.foo", "service_subnet_id"),
				),
			},
			{
				Config: testAccEksClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccTencentCloudEKSClusterExists("resource_eks_cluster.foo"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "cluster_name", "tf-ci-test"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "k8s_version", "1.12.8"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "cluster_desc", "test eks cluster updated by terraform"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.domain", "example1.org"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.servers.#", "2"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.servers.0", "10.0.0.1:82"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "dns_servers.0.servers.1", "10.0.0.1:83"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "enable_vpc_core_dns", "false"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "need_delete_cbs", "false"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "tags.hello", "world"),
					resource.TestCheckResourceAttr("resource_eks_cluster.foo", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet("resource_eks_cluster.foo", "subnet_ids.0"),
					resource.TestCheckResourceAttrSet("resource_eks_cluster.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("resource_eks_cluster.foo", "service_subnet_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_eks_cluster.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTencentCloudEKSClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("eks cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("eks cluster id is not set")
		}

		eksService := EksService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		_, has, err := eksService.DescribeEksCluster(ctx, rs.Primary.ID)

		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = eksService.DescribeEksCluster(ctx, rs.Primary.ID)
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
			return fmt.Errorf("eks cluser %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccTencentCloudEKSClusterDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	eksService := EksService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eks_cluster" {
			continue
		}
		_, has, err := eksService.DescribeEksCluster(ctx, rs.Primary.ID)

		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = eksService.DescribeEksCluster(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return err
		}

		if has {
			return fmt.Errorf("eks cluser still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

const testAccEksCluster = defaultVpcVariable + `
resource "tencentcloud_eks_cluster" "foo" {
  cluster_name = "tf-eks-test"
  k8s_version = "1.18.4"
  vpc_id = var.vpc_id
  subnet_ids = [
    var.subnet_id,
  ]
  cluster_desc = "test eks cluster created by terraform"
  service_subnet_id = var.subnet_id
  dns_servers {
    domain = "example2.org"
    servers = ["10.0.0.1:80", "10.0.0.1:81"]
  }
  enable_vpc_core_dns = true
  need_delete_cbs = true
  tags = {
    test = "tf"
  }
}
`

const testAccEksClusterUpdate = defaultVpcVariable + `
resource "tencentcloud_eks_cluster" "foo" {
  cluster_name = "tf-eks-test"
  k8s_version = "1.18.4"
  vpc_id = var.vpc_id
  subnet_ids = [
	var.subnet_id,
  ]
  cluster_desc = "test eks cluster updated by terraform"
  service_subnet_id = var.subnet_id
  dns_servers {
    domain = "example1.org"
    servers = ["10.0.0.1:82", "10.0.0.1:83"]
  }
  enable_vpc_core_dns = false
  need_delete_cbs = false
  tags = {
    test = "tf"
  }
}
`
