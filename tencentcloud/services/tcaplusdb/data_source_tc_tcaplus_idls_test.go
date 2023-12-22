package tcaplusdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTcaplusIdlsName = "data.tencentcloud_tcaplus_idls.id_test"

func TestAccTencentCloudTcaplusIdlsData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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

const testAccTencentCloudDataTcaplusIdlsBaic = tcacctest.DefaultTcaPlusData + `

resource "tencentcloud_tcaplus_idl" "test_idl" {
  cluster_id     = local.tcaplus_id
  tablegroup_id  = local.tcaplus_table_group_id
  file_name      = "tf_idl_test_guagua"
  file_type      = "PROTO"
  file_ext_type  = "proto"
  file_content   = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online_datasource {
        option(tcaplusservice.tcaplus_primary_key) = "uin,name,region";
        required int64 uin = 1;
        required string name = 2;
        required int32 region = 3;
        required int32 gamesvrid = 4;
        optional int32 logintime = 5 [default = 1];
        repeated int64 lockid = 6 [packed = true];
        optional bool is_available = 7 [default = false];
        optional pay_info_datasource pay = 8;
    }

    message pay_info_datasource {
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
