package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverBackupsName = "data.tencentcloud_sqlserver_backups.example"

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverBackups -v
func TestAccDataSourceTencentCloudSqlserverBackups(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().AddDate(0, 0, 1).In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTencentCloudDataSqlserverBackupsBasic, startTime, endTime),
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

const testAccTencentCloudDataSqlserverBackupsBasic string = `
data "tencentcloud_sqlserver_backups" "example" {
  instance_id = "mssql-qelbzgwf"
  start_time  = "%s"
  end_time    = "%s"
}
`
