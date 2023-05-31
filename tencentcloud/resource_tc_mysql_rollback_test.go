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
		},
	})
}

const testAccMysqlRollbackVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlRollback = testAccMysqlRollbackVar + `

resource "tencentcloud_mysql_rollback" "rollback" {
	instance_id = var.instance_id
	strategy = "full"
	rollback_time = "2023-05-31 23:13:35"
	databases {
	  database_name = "tf_ci_test_bak"
	  new_database_name = "tf_ci_test_bak_5"
	}
	tables {
	  database = "tf_ci_test_bak"
	  table {
		table_name = "test"
		new_table_name = "test_bak"
	  }
	}
}

`
