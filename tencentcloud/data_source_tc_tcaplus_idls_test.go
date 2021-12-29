package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusIdlsName = "data.tencentcloud_tcaplus_idls.id_test"

func TestAccTencentCloudDataTcaplusIdls(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusIdlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusIdlsBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusIdlExists("tencentcloud_tcaplus_idl.test_idl"),
					resource.TestCheckResourceAttrSet(testDataTcaplusIdlsName, "cluster_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusIdlsName, "list.#"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusIdlsBaic = `
variable "availability_zone" {
default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_tablegroup" "test_group" {
  cluster_id         = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_name    = "tf_test_group_name_guagua"
}

resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_data_guagua"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
resource "tencentcloud_tcaplus_idl" "test_idl" {
  cluster_id     = tencentcloud_tcaplus_cluster.test_cluster.id
  tablegroup_id  = tencentcloud_tcaplus_tablegroup.test_group.id
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
   cluster_id     = tencentcloud_tcaplus_idl.test_idl.cluster_id
}
`
