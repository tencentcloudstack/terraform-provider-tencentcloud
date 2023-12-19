package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMysqlBackupListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlBackupListConfig(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_backup_list.test"),
				),
			},
		},
	})
}

func testAccDataSourceMysqlBackupListConfig() string {
	return tcacctest.CommonPresetMysql + `
data "tencentcloud_mysql_backup_list" "test" {
  mysql_id = local.mysql_id
}`
}
