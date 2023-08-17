package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverQueryXeventDataSource_basic -v
func TestAccTencentCloudSqlserverQueryXeventDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().AddDate(0, 0, 1).In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSqlserverQueryXeventDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_query_xevent.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.example", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.example", "event_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.example", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.example", "end_time"),
				),
			},
		},
	})
}

const testAccSqlserverQueryXeventDataSource = `
data "tencentcloud_sqlserver_query_xevent" "example" {
  instance_id = "mssql-gyg9xycl"
  event_type  = "blocked"
  start_time  = "%s"
  end_time    = "%s"
}
`
