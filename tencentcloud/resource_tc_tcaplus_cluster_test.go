package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTcaplusClusterResourceName = "tencentcloud_tcaplus_cluster"
var testTcaplusClusterResourceKey = testTcaplusClusterResourceName + ".test_cluster"

func TestAccTencentCloudTcaplusClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusClusterExists(testTcaplusClusterResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_port"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "cluster_name", "tf_te1_guagua"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "old_password_expire_last", "3600"),
				),
			},
			{
				ResourceName:            testTcaplusClusterResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"old_password_expire_last", "password"},
			},

			{
				Config: testAccTcaplusClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusClusterExists(testTcaplusClusterResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_port"),

					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "cluster_name", "tf_te1_guagua_2"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "password", "aQQ2345677888"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "old_password_expire_last", "300"),
				),
			},
		},
	})
}

func testAccCheckTcaplusClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusClusterResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete tcaplus cluster %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTcaplusClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("tcaplus cluster %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccTcaplusClusterBaic = `
variable "availability_zone" {
default = "ap-shanghai-2"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}
`

const testAccTcaplusCluster string = testAccTcaplusClusterBaic + `
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_te1_guagua"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
`
const testAccTcaplusClusterUpdate string = testAccTcaplusClusterBaic + `
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_te1_guagua_2"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "aQQ2345677888"
  old_password_expire_last = 300
}
`
