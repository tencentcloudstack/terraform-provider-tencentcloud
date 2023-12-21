package dts_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDtsMigrateDbInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateDbInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dts_migrate_db_instances.migrate_db_instances")),
			},
		},
	})
}

const testAccDtsMigrateDbInstancesDataSource = tcacctest.CommonPresetMysql + `

data "tencentcloud_dts_migrate_db_instances" "migrate_db_instances" {
  database_type = "mysql"
  migrate_role = "src"
  instance_id = local.mysql_id
  instance_name = "test"
  account_mode = "self"
  tmp_secret_id = "AKIDvBDyVmna9TadcS4YzfBZmkU5TbX12345"
  tmp_secret_key = "ZswjGWWHm24qMeiX6QUJsELDpC12345"
  tmp_token = "JOqqCPVuWdNZvlVDLxxx"
}

`
