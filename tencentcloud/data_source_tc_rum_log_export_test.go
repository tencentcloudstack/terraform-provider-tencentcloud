package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixRumLogExportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumLogExportDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_log_export.log_export")),
			},
		},
	})
}

const testAccRumLogExportDataSource = `

data "tencentcloud_rum_group_log" "group_log" {
  order_by    = "desc"
  start_time  = 1696216110
  query       = "type:\"log\""
  end_time    = 1696820910
  project_id  = 120000
  group_field = "level"
}

`
