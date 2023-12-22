package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainSqlFiltersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbbrainSqlFilters(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_sql_filters.sql_filters"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_sql_filters.sql_filters", "list.#"),
				),
			},
		},
	})
}

func testAccDataSourceDbbrainSqlFilters() string {
	return fmt.Sprintf(`%s

resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = local.mysql_id
  session_token {
    user = "keep_dbbrain"
	password = "Test@123456#"
  }
  sql_type = "SELECT"
  filter_key = "test"
  max_concurrency = 10
  duration = 3600
}

data "tencentcloud_dbbrain_sql_filters" "sql_filters" {
  instance_id = local.mysql_id
  filter_ids = [tencentcloud_dbbrain_sql_filter.sql_filter.filter_id]
  }
  `, tcacctest.CommonPresetMysql)
}
