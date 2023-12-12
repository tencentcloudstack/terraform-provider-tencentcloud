package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverDescHaLogDataSource_basic -v
func TestAccTencentCloudSqlserverDescHaLogDataSource_basic(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSqlserverDescHaLogDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_desc_ha_log.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_desc_ha_log.example", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_desc_ha_log.example", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_desc_ha_log.example", "end_time"),
				),
			},
		},
	})
}

const testAccSqlserverDescHaLogDataSource = `
data "tencentcloud_sqlserver_desc_ha_log" "example" {
  instance_id = "mssql-gyg9xycl"
  start_time  = "%s"
  end_time    = "%s"
  switch_type = 1
}
`
