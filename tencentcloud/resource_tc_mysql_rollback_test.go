package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlRollbackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRollback,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_rollback.rollback", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_rollback.rollback",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlRollback = `

resource "tencentcloud_mysql_rollback" "rollback" {
  instances {
		instance_id = "cdb_xxx"
		strategy = ""
		rollback_time = ""
		databases {
			database_name = ""
			new_database_name = ""
		}
		tables {
			database = ""
			table {
				table_name = ""
				new_table_name = ""
			}
		}

  }
}

`
