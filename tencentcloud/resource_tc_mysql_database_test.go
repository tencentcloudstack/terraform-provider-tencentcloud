package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlDatabaseResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDatabase,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_database.database", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_database.database",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlDatabase = `

resource "tencentcloud_mysql_database" "database" {
  instance_id        = "cdb-i9xfdf7z"
  db_name            = "for_tf_test"
  character_set_name = "utf8"
}

`
