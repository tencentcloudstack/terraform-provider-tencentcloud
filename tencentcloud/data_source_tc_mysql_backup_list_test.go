package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMysqlBackupListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlBackupListConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_backup_list.test"),
				),
			},
		},
	})
}

func testAccDataSourceMysqlBackupListConfig() string {
	return CommonPresetMysql + `
data "tencentcloud_mysql_backup_list" "test" {
  mysql_id = local.mysql_id
}`
}
