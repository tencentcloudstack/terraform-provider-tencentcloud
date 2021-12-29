package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTcaplusGroupResourceName = "tencentcloud_tcaplus_tablegroup"
var testTcaplusGroupResourceNameResourceKey = testTcaplusGroupResourceName + ".test_group"

func TestAccTencentCloudTcaplusGroupResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusGroupExists(testTcaplusGroupResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusGroupResourceNameResourceKey, "total_size"),
					resource.TestCheckResourceAttrSet(testTcaplusGroupResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testTcaplusGroupResourceNameResourceKey, "tablegroup_name", "tf_test_group_name_guagua"),
					resource.TestCheckResourceAttr(testTcaplusGroupResourceNameResourceKey, "table_count", "0"),
				),
			},
			{
				Config: testAccTcaplusGroupUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusGroupExists(testTcaplusGroupResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusGroupResourceNameResourceKey, "total_size"),
					resource.TestCheckResourceAttrSet(testTcaplusGroupResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testTcaplusGroupResourceNameResourceKey, "tablegroup_name", "tf_test_group_name_guagua_2"),
					resource.TestCheckResourceAttr(testTcaplusGroupResourceNameResourceKey, "table_count", "0"),
				),
			},
		},
	})
}
func testAccCheckTcaplusGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusGroupResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeGroup(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)

		if err != nil {
			_, has, err = service.DescribeGroup(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}

		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("delete tcaplus group %s fail, still on server", rs.Primary.ID)
	}
	return nil
}

func testAccCheckTcaplusGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeGroup(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)

		if err != nil {
			_, has, err = service.DescribeGroup(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}

		if err != nil {
			return err
		}
		if has {
			return nil
		}
		return fmt.Errorf("tcaplus group %s not found on server", rs.Primary.ID)
	}
}

const testAccTcaplusGroupBasic = `variable "availability_zone" {
  default = "ap-guangzhou-3"
}
data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_guagua_group"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}`

const testAccTcaplusGroup = testAccTcaplusGroupBasic + `
resource "tencentcloud_tcaplus_tablegroup" "test_group" {
  cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name    = "tf_test_group_name_guagua"
}
`
const testAccTcaplusGroupUpdate = testAccTcaplusGroupBasic + `
resource "tencentcloud_tcaplus_tablegroup" "test_group" {
  cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name    = "tf_test_group_name_guagua_2"
}
`
