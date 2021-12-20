package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTcaplusTableResourceName = "tencentcloud_tcaplus_table"
var testTcaplusTableResourceNameResourceKey = testTcaplusTableResourceName + ".test_table"

func TestAccTencentCloudTcaplusTableResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusTable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusTableExists(testTcaplusTableResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "cluster_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "tablegroup_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "status"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "table_size"),

					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_name", "tb_online_guagua"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_type", "GENERIC"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "description", "test"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_read_cu", "1000"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_write_cu", "20"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_volume", "1"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "error", ""),
				),
			},
			{
				Config: testAccTcaplusTableUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusTableExists(testTcaplusTableResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "cluster_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "tablegroup_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "status"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "table_size"),

					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_name", "tb_online_guagua"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_type", "GENERIC"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "description", "test_desc"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_read_cu", "1000"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_write_cu", "20"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_volume", "1"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "error", ""),
				),
			},
		},
	})
}
func testAccCheckTcaplusTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusTableResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeTable(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeTable(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
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

func testAccCheckTcaplusTableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeTable(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeTable(ctx, rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
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

const testAccTcaplusTableBasic = `variable "availability_zone" {
  default = "ap-shanghai-2"
}
data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_g_table"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
resource "tencentcloud_tcaplus_tablegroup" "group" {
  cluster_id           = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name      = "tf_test_group_name"
}
resource "tencentcloud_tcaplus_idl" "test_idl" {
  cluster_id     = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id  = tencentcloud_tcaplus_tablegroup.group.id
  file_name      = "tf_idl_test_guagua"
  file_type      = "PROTO"
  file_ext_type  = "proto"
  file_content   = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online_guagua {
        option(tcaplusservice.tcaplus_primary_key) = "uin,name,region";
        required int64 uin = 1;
        required string name = 2;
        required int32 region = 3;
        required int32 gamesvrid = 4;
        optional int32 logintime = 5 [default = 1];
        repeated int64 lockid = 6 [packed = true];
        optional bool is_available = 7 [default = false];
        optional pay_info pay = 8;
    }

    message pay_info {
        required int64 pay_id = 1;
        optional uint64 total_money = 2;
        optional uint64 pay_times = 3;
        optional pay_auth_info auth = 4;
        message pay_auth_info {
            required string pay_keys = 1;
            optional int64 update_time = 2;
        }
    }
    EOF
}
resource "tencentcloud_tcaplus_tablegroup" "test_group" {
  cluster_id      = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name = "tf_test_group_name_guagua"
}

resource "tencentcloud_tcaplus_idl" "test_idl_2" {
  cluster_id     = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id  = tencentcloud_tcaplus_tablegroup.test_group.id
  file_name      = "tf_idl_test_guagua_2"
  file_type      = "PROTO"
  file_ext_type  = "proto"
  file_content   = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online_guagua {
        option(tcaplusservice.tcaplus_primary_key) = "uin,name,region";
        required int64 uin = 1;
        required string name = 2;
        required int32 region = 3;
        required int32 gamesvrid = 4;
        optional int32 logintime = 5 [default = 1];
        repeated int64 lockid = 6 [packed = true];
        optional bool is_available = 7 [default = false];
        optional pay_info pay = 8;
    }

    message pay_info {
        required int64 pay_id = 1;
        optional uint64 total_money = 2;
        optional uint64 pay_times = 3;
        optional pay_auth_info auth = 4;
        message pay_auth_info {
            required string pay_keys = 1;
            optional int64 update_time = 2;
        }
    }
    EOF
}
`
const testAccTcaplusTable = testAccTcaplusTableBasic + `
resource "tencentcloud_tcaplus_table" "test_table" {
  cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id      = tencentcloud_tcaplus_tablegroup.test_group.id
  table_name         = "tb_online_guagua"
  table_type         = "GENERIC"
  description        = "test"
  idl_id             = tencentcloud_tcaplus_idl.test_idl.id
  table_idl_type     = "PROTO"
  reserved_read_cu   = 1000
  reserved_write_cu  = 20
  reserved_volume    = 1
}
`
const testAccTcaplusTableUpdate = testAccTcaplusTableBasic + `
resource "tencentcloud_tcaplus_table" "test_table" {
  cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id      = tencentcloud_tcaplus_tablegroup.test_group.id
  table_name         = "tb_online_guagua"
  table_type         = "GENERIC"
  description        = "test_desc"
  idl_id             = tencentcloud_tcaplus_idl.test_idl_2.id
  table_idl_type     = "PROTO"
  reserved_read_cu  = 1000
  reserved_write_cu = 20
  reserved_volume    = 1
}
`
