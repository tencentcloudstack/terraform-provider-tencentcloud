package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainNoPrimaryKeyTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	// loc, _ := time.LoadLocation("Asia/Chongqing")
	// queryDate := time.Now().AddDate(0, 0, -20).In(loc).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainNoPrimaryKeyTablesDataSource, "2023-06-13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "date", "2023-06-13"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "product", "mysql"),
					// return
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_tables.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_tables.0.table_schema"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_tables.0.table_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_tables.0.engine"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_tables.0.table_rows"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_tables.0.total_length"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "timestamp"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables", "no_primary_key_table_count_diff"),
				),
			},
		},
	})
}

const testAccDbbrainNoPrimaryKeyTablesDataSource = CommonPresetMysql + `

data "tencentcloud_dbbrain_no_primary_key_tables" "no_primary_key_tables" {
  instance_id = local.mysql_id
  date = "%s"
  product = "mysql"
}

`
