package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbRollbackClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbRollbackCluster,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_rollback_cluster.rollback_cluster", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_rollback_cluster.rollback_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbRollbackCluster = `

resource "tencentcloud_cynosdb_rollback_cluster" "rollback_cluster" {
  cluster_id = "cynosdbmysql-xxxxxxxx"
  rollback_strategy = "timeRollback"
  rollback_id = 1
  expect_time = "	2022-01-20 00:00:00"
  expect_time_thresh = 1
  rollback_databases {
		old_database = "old_db_1"
		new_database = "new_db_1"

  }
  rollback_tables {
		database = "old_db_1"
		tables {
			old_table = "old_tbl_1"
			new_table = "new_tbl_1"
		}

  }
  rollback_mode = "full"
}

`
