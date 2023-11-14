package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbExportInstanceSlowQueriesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbExportInstanceSlowQueries,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_slow_queries.export_instance_slow_queries", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_export_instance_slow_queries.export_instance_slow_queries",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbExportInstanceSlowQueries = `

resource "tencentcloud_cynosdb_export_instance_slow_queries" "export_instance_slow_queries" {
  instance_id = "cynosdbmysql-ins-123"
  start_time = "2022-01-01 12:00:00"
  end_time = "2022-01-01 14:00:00"
  username = "root"
  host = "10.10.10.10"
  database = "db1"
  file_type = "csv"
}

`
