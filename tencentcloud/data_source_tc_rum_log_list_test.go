package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixRumLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumLogListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_log_list.log_list")),
			},
		},
	})
}

const testAccRumLogListDataSource = `

data "tencentcloud_rum_log_list" "log_list" {
  order_by   = "desc"
  start_time = 1696216110
  query      = "id:123 AND type:\"log\""
  end_time   = 1696820910
  project_id = 120000
}

`
