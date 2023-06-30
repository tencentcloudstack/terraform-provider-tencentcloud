package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverQueryXeventDataSource_basic -v
func TestAccTencentCloudSqlserverQueryXeventDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverQueryXeventDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_query_xevent.query_xevent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.query_xevent", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.query_xevent", "event_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.query_xevent", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_query_xevent.query_xevent", "end_time"),
				),
			},
		},
	})
}

const testAccSqlserverQueryXeventDataSource = `
data "tencentcloud_sqlserver_query_xevent" "query_xevent" {
  instance_id = "mssql-gyg9xycl"
  event_type  = "blocked"
  start_time  = "2023-06-27 00:00:00"
  end_time    = "2023-07-01 00:00:00"
}
`
