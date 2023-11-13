package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_query_xevent.query_xevent")),
			},
		},
	})
}

const testAccSqlserverQueryXeventDataSource = `

data "tencentcloud_sqlserver_query_xevent" "query_xevent" {
  instance_id = ""
  event_type = ""
  start_time = ""
  end_time = ""
  }

`
