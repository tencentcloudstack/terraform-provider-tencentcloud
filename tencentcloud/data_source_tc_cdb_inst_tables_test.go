package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbInstTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbInstTablesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_inst_tables.inst_tables")),
			},
		},
	})
}

const testAccCdbInstTablesDataSource = `

data "tencentcloud_cdb_inst_tables" "inst_tables" {
  instance_id = ""
  database = ""
  table_regexp = ""
  }

`
