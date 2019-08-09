package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMysqlBackupList_basic(t *testing.T) {
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
	return MysqlInstanceCommonTestCase + `
data "tencentcloud_mysql_backup_list" "test" {
  mysql_id = "${tencentcloud_mysql_instance.default.id}"
}`
}
