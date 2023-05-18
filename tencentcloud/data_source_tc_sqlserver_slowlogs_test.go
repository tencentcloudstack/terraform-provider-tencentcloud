package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverSlowlogsDataSource_basic -v
func TestAccTencentCloudSqlserverSlowlogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverSlowlogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_slowlogs.slowlogs")),
			},
		},
	})
}

const testAccSqlserverSlowlogsDataSource = `
data "tencentcloud_sqlserver_slowlogs" "slowlogs" {
  instance_id = "mssql-qelbzgwf"
  start_time = "2020-05-01 00:00:00"
  end_time = "2023-05-18 00:00:00"
}
`
