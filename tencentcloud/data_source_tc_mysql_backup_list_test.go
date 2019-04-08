package tencentcloud

import (
	"fmt"
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
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.finish_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.backup_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.backup_model"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.intranet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.internet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.test", "list.0.creator"),
				),
			},
			{
				Config: testAccDataSourceMysqlBackupListConfigFull(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_backup_list.testFull"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.finish_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.backup_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.backup_model"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.intranet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.internet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_list.testFull", "list.0.creator"),
				),
			},
		},
	})
}

/*
	you should modify to your mysql_id before testing.
*/
const mysqlIdfor_TestAccDataSourceMysqlBackupListConfig = "cdb-ia8zhj0t"

func testAccDataSourceMysqlBackupListConfig() string {
	return fmt.Sprintf(`
data "tencentcloud_mysql_backup_list" "test" {
		mysql_id = "%s"
}`, mysqlIdfor_TestAccDataSourceMysqlBackupListConfig)
}

func testAccDataSourceMysqlBackupListConfigFull() string {

	return fmt.Sprintf(`
data "tencentcloud_mysql_backup_list" "testFull" {
		mysql_id = "%s"
		max_number = 100
		result_output_file ="/tmp/backup_list"
}
`, mysqlIdfor_TestAccDataSourceMysqlBackupListConfig)

}
