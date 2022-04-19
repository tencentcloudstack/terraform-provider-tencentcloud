package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverBackupsName = "data.tencentcloud_sqlserver_backups.test"

var now = time.Now().Format("2006-01-02 15:04:05")

func TestAccDataSourceTencentCloudSqlserverBackups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverBackupsBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSqlserverBackupsName, "end_time", now),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.start_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.end_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.file_name"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.size"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.strategy"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.trigger_model"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.intranet_url"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.internet_url"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverBackupsName, "list.0.db_list.0"),
				),
			},
		},
	})
}

func testAccTencentCloudDataSqlserverBackupsBasic() string {
	return fmt.Sprintf(`
%s
data "tencentcloud_sqlserver_backups" "test"{
	instance_id = local.sqlserver_id
	start_time = "2020-06-17 00:00:00"
	end_time = "%s"
}
`, CommonPresetSQLServer, now)
}
