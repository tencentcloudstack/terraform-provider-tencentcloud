package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusTablesName = "data.tencentcloud_tcaplus_tables.id_test"

func TestAccTencentCloudTcaplusTablesData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusTablesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "cluster_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "table_name"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.#"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.tablegroup_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.table_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.table_name"),
					resource.TestCheckResourceAttr(testDataTcaplusTablesName, "list.0.table_type", "GENERIC"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.description"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.idl_id"),
					resource.TestCheckResourceAttr(testDataTcaplusTablesName, "list.0.table_idl_type", "PROTO"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.reserved_read_cu"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.reserved_write_cu"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.reserved_volume"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.status"),
					resource.TestCheckResourceAttrSet(testDataTcaplusTablesName, "list.0.table_size"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusTablesBasic = defaultTcaPlusData + `
data "tencentcloud_tcaplus_tables" "id_test" {
  cluster_id   = local.tcaplus_id
  table_name     =  local.tcaplus_table
}
`
