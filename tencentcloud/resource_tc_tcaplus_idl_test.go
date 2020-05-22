package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTcaplusIdlResourceName = "tencentcloud_tcaplus_idl"
var testTcaplusIdlResourceNameResourceKey = testTcaplusIdlResourceName + ".test_idl"

func TestAccTencentCloudTcaplusIdlResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusIdlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusIdl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusIdlExists(testTcaplusIdlResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusIdlResourceNameResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTcaplusIdlResourceNameResourceKey, "file_name", "tf_idl_test_guagua"),
					resource.TestCheckResourceAttr(testTcaplusIdlResourceNameResourceKey, "file_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusIdlResourceNameResourceKey, "file_ext_type", "proto"),

					resource.TestCheckResourceAttr(testTcaplusIdlResourceNameResourceKey, "table_infos.#", "1"),

					resource.TestCheckResourceAttrSet(testTcaplusIdlResourceNameResourceKey, "table_infos.0.key_fields"),
					resource.TestCheckResourceAttrSet(testTcaplusIdlResourceNameResourceKey, "table_infos.0.sum_key_field_size"),
					resource.TestCheckResourceAttrSet(testTcaplusIdlResourceNameResourceKey, "table_infos.0.value_fields"),
					resource.TestCheckResourceAttrSet(testTcaplusIdlResourceNameResourceKey, "table_infos.0.sum_value_field_size"),
					resource.TestCheckResourceAttrSet(testTcaplusIdlResourceNameResourceKey, "table_infos.0.index_key_set"),
					resource.TestCheckResourceAttr(testTcaplusIdlResourceNameResourceKey, "table_infos.0.table_name", "tb_online_guagua"),
					resource.TestCheckResourceAttr(testTcaplusIdlResourceNameResourceKey, "table_infos.0.error", ""),
				),
			},
		},
	})
}
func testAccCheckTcaplusIdlDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusIdlResourceName {
			continue
		}
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		var tcaplusIdlId TcaplusIdlId

		if err := json.Unmarshal([]byte(rs.Primary.ID), &tcaplusIdlId); err != nil {
			return fmt.Errorf("idl id is broken,%s", err.Error())
		}
		parseTableInfos, err := service.DesOldIdlFiles(ctx, tcaplusIdlId)
		if err != nil {
			parseTableInfos, err = service.DesOldIdlFiles(ctx, tcaplusIdlId)
		}
		if err != nil {
			return err
		}
		if len(parseTableInfos) == 0 {
			return nil
		}
		return fmt.Errorf("delete tcaplus idl %s fail, still on server", rs.Primary.ID)
	}
	return nil
}

func testAccCheckTcaplusIdlExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		var tcaplusIdlId TcaplusIdlId

		if err := json.Unmarshal([]byte(rs.Primary.ID), &tcaplusIdlId); err != nil {
			return fmt.Errorf("idl id is broken,%s", err.Error())
		}
		parseTableInfos, err := service.DesOldIdlFiles(ctx, tcaplusIdlId)
		if err != nil {
			parseTableInfos, err = service.DesOldIdlFiles(ctx, tcaplusIdlId)
		}
		if err != nil {
			return err
		}
		if len(parseTableInfos) != 0 {
			return nil
		}
		return fmt.Errorf("tcaplus idl %s not found on server", rs.Primary.ID)
	}
}

const testAccTcaplusIdlBasic = `variable "availability_zone" {
  default = "ap-shanghai-2"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_test_cluster"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
resource "tencentcloud_tcaplus_group" "group" {
  cluster_id     = tencentcloud_tcaplus_cluster.test_cluster.id
  group_name     = "tf_test_group_name"
}
`

const testAccTcaplusIdl = testAccTcaplusIdlBasic + `
resource "tencentcloud_tcaplus_idl" "test_idl" {
  cluster_id     = tencentcloud_tcaplus_cluster.test_cluster.id
  group_id       = tencentcloud_tcaplus_group.group.id
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
`
