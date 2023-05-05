package tencentcloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_tdcpg_cluster", &resource.Sweeper{
		Name: "tencentcloud_tdcpg_cluster",
		F:    testSweepTdcpgCluster,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_tdcpg_cluster
func testSweepTdcpgCluster(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	tdcpgService := TdcpgService{client: cli.(*TencentCloudClient).apiV3Conn}

	clusters, err := tdcpgService.DescribeTdcpgClustersByFilter(ctx, nil)
	if err != nil {
		return err
	}
	if clusters == nil {
		return fmt.Errorf("No any tdcpg clusters exists.")
	}

	// delete all cluster with specified prefix
	for _, v := range clusters {
		delId := v.ClusterId
		delName := v.ClusterName
		status := *v.Status

		if status == "deleted" {
			continue
		}
		if strings.HasPrefix(*delName, defaultTdcpgTestNamePrefix) {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := tdcpgService.DeleteTdcpgClusterById(ctx, delId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] delete tdcpg cluster %s failed. reason:[%s]", *delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudTdcpgClusterResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdcpgClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTdcpgCluster_postpaid, defaultTdcpgZone, defaultTdcpgTestNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdcpgClusterExists("tencentcloud_tdcpg_cluster.cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "zone", defaultTdcpgZone),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "cpu", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "memory", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "pay_mode", "POSTPAID_BY_HOUR"),
					resource.TestMatchResourceAttr("tencentcloud_tdcpg_cluster.cluster", "cluster_name", regexp.MustCompile(defaultTdcpgTestNamePrefix)),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "db_version", "10.17"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "instance_count", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "period"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "storage"),
					resource.TestCheckResourceAttr("tencentcloud_tdcpg_cluster.cluster", "project_id", "0"),
				),
			},
			{
				ResourceName:            "tencentcloud_tdcpg_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"master_user_password", "period"},
			},
		},
	})
}

func testAccCheckTdcpgClusterDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tdcpgService := TdcpgService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdcpg_cluster" {
			continue
		}

		ret, err := tdcpgService.DescribeTdcpgCluster(ctx, &rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(ret.ClusterSet) > 0 {
			status := *ret.ClusterSet[0].Status
			if status == "deleting" || status == "deleted" || status == "isolated" || status == "isolating" {
				return nil
			}
			return fmt.Errorf("tdcpg cluster still exist, clusterId: %v, status: %v", rs.Primary.ID, status)
		}
	}
	return nil
}

func testAccCheckTdcpgClusterExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("tdcpg cluster instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("tdcpg cluster instance id is not set")
		}

		tdcpgService := TdcpgService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ret, err := tdcpgService.DescribeTdcpgCluster(ctx, &rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(ret.ClusterSet) == 0 {
			return fmt.Errorf("tdcpg cluster instance not found, clusterId: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdcpg_vpc_config = `
data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}
	
data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}
	
locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	#sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}
`

const testAccTdcpgCluster_postpaid = testAccTdcpg_vpc_config + `

resource "tencentcloud_tdcpg_cluster" "cluster" {
  zone = "%s"
  master_user_password = "===Password123==="
  cpu = 1
  memory = 1
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  pay_mode = "POSTPAID_BY_HOUR"
  cluster_name = "%scluster"
  db_version = "10.17"
  instance_count = 1
  period = 1
  project_id = 0
}

`
