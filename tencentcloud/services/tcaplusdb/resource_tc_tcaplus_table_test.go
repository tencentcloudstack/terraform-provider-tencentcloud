package tcaplusdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcaplusdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcaplusdb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testTcaplusTableResourceName = "tencentcloud_tcaplus_table"
var testTcaplusTableResourceNameResourceKey = testTcaplusTableResourceName + ".test_table"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_tcaplus_table
	resource.AddTestSweepers("tencentcloud_tcaplus_table", &resource.Sweeper{
		Name: "tencentcloud_tcaplus_table",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svctcaplusdb.NewTcaplusService(client)

			clusters, err := service.DescribeClusters(ctx, "", tcacctest.DefaultTcaPlusClusterName)

			if err != nil {
				return err
			}

			if len(clusters) == 0 {
				return fmt.Errorf("no cluster named %s", tcacctest.DefaultTcaPlusClusterName)
			}

			clusterId := *clusters[0].ClusterId

			tables, _ := service.DescribeTables(ctx, clusterId, "", "", "")

			for _, table := range tables {
				name := *table.TableName
				// legacy bad argument
				gId := clusterId + ":" + *table.TableGroupId
				insId := *table.TableInstanceId
				created := time.Time{}

				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}

				taskId, err := service.DeleteTable(ctx, clusterId, gId, insId, name)
				if err != nil {
					continue
				}

				err = resource.Retry(tccommon.ReadRetryTimeout*3, func() *resource.RetryError {
					info, has, err := service.DescribeTask(ctx, clusterId, taskId)
					if err != nil {
						return tccommon.RetryError(err)
					}
					if !has {
						return nil
					}
					if *info.Progress < 100 {
						return resource.RetryableError(fmt.Errorf("running delete task %s, table: %s -> %s", taskId, clusterId, name))
					}
					return nil
				})

				if err != nil {
					continue
				}

				_, err = service.DeleteTable(ctx, clusterId, gId, insId, name)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudTcaplusTableResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svctcaplusdb.NewTcaplusService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svctcaplusdb.NewTcaplusService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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

const testAccTcaplusTableBasic = tcacctest.DefaultTcaPlusData + `
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

resource "tencentcloud_tcaplus_idl" "test_idl_2" {
  cluster_id     = local.tcaplus_id
  tablegroup_id  = local.tcaplus_table_group_id
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
  cluster_id         = local.tcaplus_id
  tablegroup_id      = local.tcaplus_table_group_id
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
  cluster_id         = local.tcaplus_id
  tablegroup_id      = local.tcaplus_table_group_id
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
