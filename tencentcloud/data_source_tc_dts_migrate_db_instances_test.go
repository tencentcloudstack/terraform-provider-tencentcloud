package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateDbInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateDbInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dts_migrate_db_instances.migrate_db_instances")),
			},
		},
	})
}

const testAccDtsMigrateDbInstancesDataSource = `

data "tencentcloud_dts_migrate_db_instances" "migrate_db_instances" {
  database_type = "mysql"
  migrate_role = "src"
  instance_id = "cdb-ffulb2sg"
  instance_name = "cdb_test"
  limit = 10
  offset = 10
  account_mode = "self"
  tmp_secret_id = ""
  tmp_secret_key = ""
  tmp_token = ""
      }

`
