package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusIdlsName = "data.tencentcloud_tcaplus_idls.id_test"

func TestAccTencentCloudDataTcaplusIdls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusIdlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusIdlsBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusIdlExists("tencentcloud_tcaplus_idl.test_idl"),
					resource.TestCheckResourceAttrSet(testDataTcaplusIdlsName, "app_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusIdlsName, "list.#"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusIdlsBaic = `
variable "availability_zone" {
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
  app_name                 = "tf_tcaplus_data_guagua"
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
data "tencentcloud_tcaplus_idls" "id_test" {
   app_id     = tencentcloud_tcaplus_idl.test_idl.app_id
}
`
