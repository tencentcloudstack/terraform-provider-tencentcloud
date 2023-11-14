package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRollbackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRollback,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_rollback.rollback", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_rollback.rollback",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRollback = `

resource "tencentcloud_cdb_rollback" "rollback" {
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
