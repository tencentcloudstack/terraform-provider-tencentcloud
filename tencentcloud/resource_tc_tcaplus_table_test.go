package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

var testTcaplusTableResourceName = "tencentcloud_tcaplus_table"
var testTcaplusTableResourceNameResourceKey = testTcaplusTableResourceName + ".test_table"

func TestAccTencentCloudTcaplusTableResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusTable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusTableExists(testTcaplusTableResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "app_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "zone_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "status"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "table_size"),

					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_name", "tb_online_guagua"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_type", "GENERIC"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "description", "test"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_read_qps", "1000"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_write_qps", "20"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_volume", "1"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "error", ""),
				),
			},
			{
				Config: testAccTcaplusTableUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusTableExists(testTcaplusTableResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "app_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "zone_id"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "status"),
					resource.TestCheckResourceAttrSet(testTcaplusTableResourceNameResourceKey, "table_size"),

					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_name", "tb_online_guagua"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_type", "GENERIC"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "description", "test_desc"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "table_idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_read_qps", "1000"),
					resource.TestCheckResourceAttr(testTcaplusTableResourceNameResourceKey, "reserved_write_qps", "20"),
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
		ctx := context.WithValue(context.TODO(), "logId", logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeTable(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeTable(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)
		}

		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("delete tcaplus zone %s fail, still on server", rs.Primary.ID)
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
		ctx := context.WithValue(context.TODO(), "logId", logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeTable(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeTable(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		}
		return fmt.Errorf("tcaplus zone %s not found on server", rs.Primary.ID)
	}
}

const testAccTcaplusTableBasic = `variable "availability_zone" {
  default = "ap-shanghai-2"
}
variable "instance_name" {
  default = "` + defaultInsName + `"
}
variable "vpc_cidr" {
  default = "` + defaultVpcCidr + `"
}
variable "subnet_cidr" {
  default = "` + defaultSubnetCidr + `"
}

resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}
resource "tencentcloud_subnet" "subnet" {
  name              = var.instance_name
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr
  is_multicast      = false
}
resource "tencentcloud_tcaplus_application" "test_app" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_g_table"
  vpc_id                   = tencentcloud_vpc.foo.id
  subnet_id                = tencentcloud_subnet.subnet.id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
resource "tencentcloud_tcaplus_idl" "test_idl" {
  app_id = tencentcloud_tcaplus_application.test_app.id
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
resource "tencentcloud_tcaplus_zone" "test_zone" {
  app_id    = tencentcloud_tcaplus_application.test_app.id
  zone_name = "tf_test_zone_name_guagua"
}

resource "tencentcloud_tcaplus_idl" "test_idl_2" {
  app_id = tencentcloud_tcaplus_application.test_app.id
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
  app_id             = tencentcloud_tcaplus_application.test_app.id
  zone_id            = tencentcloud_tcaplus_zone.test_zone.id
  table_name         = "tb_online_guagua"
  table_type         = "GENERIC"
  description        = "test"
  idl_id             = tencentcloud_tcaplus_idl.test_idl.id
  table_idl_type     = "PROTO"
  reserved_read_qps  = 1000
  reserved_write_qps = 20
  reserved_volume    = 1
}
`
const testAccTcaplusTableUpdate = testAccTcaplusTableBasic + `
resource "tencentcloud_tcaplus_table" "test_table" {
  app_id             = tencentcloud_tcaplus_application.test_app.id
  zone_id            = tencentcloud_tcaplus_zone.test_zone.id
  table_name         = "tb_online_guagua"
  table_type         = "GENERIC"
  description        = "test_desc"
  idl_id             = tencentcloud_tcaplus_idl.test_idl_2.id
  table_idl_type     = "PROTO"
  reserved_read_qps  = 1000
  reserved_write_qps = 20
  reserved_volume    = 1
}
`
