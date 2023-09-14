package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudRumLogExportDataSource_basic(t *testing.T) {
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

data "tencentcloud_rum_log_export" "log_export" {
  name = "log"
  start_time = "1692594840000"
  query = "id:123 AND type:&quot;log&quot;"
  end_time = "1692609240000"
  project_id = 1
  fields = 
  }

`
