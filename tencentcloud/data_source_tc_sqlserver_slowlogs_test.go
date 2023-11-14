package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
  instance_id = "mssql-j8kv137v"
  start_time = ""
  end_time = ""
  }

`
