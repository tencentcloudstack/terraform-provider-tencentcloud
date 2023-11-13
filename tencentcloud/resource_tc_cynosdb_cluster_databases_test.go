package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterDatabasesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterDatabases,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_databases.cluster_databases", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_databases.cluster_databases",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterDatabases = `

resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  db_name = "test"
  character_set = "utf8"
  collate_rule = " utf8_general_ci "
  user_host_privileges {
		db_user_name = ""
		db_host = ""
		db_privilege = ""

  }
  description = "test"
}

`
