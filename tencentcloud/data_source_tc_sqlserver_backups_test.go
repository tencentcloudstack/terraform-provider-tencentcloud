package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverBackupsName = "data.tencentcloud_sqlserver_backups.test"

func TestAccDataSourceTencentCloudSqlserverBackups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverBackupsBasic,
				Check: resource.ComposeTestCheckFunc(
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

const testAccTencentCloudDataSqlserverBackupsBasic = testAccSqlserverDB_basic + `
data "tencentcloud_sqlserver_backups" "test"{
	instance_id = tencentcloud_sqlserver_instance.test.id
	start_time = "2020-06-17 00:00:00"
	end_time = "2020-06-22 00:00:00"
}
`
