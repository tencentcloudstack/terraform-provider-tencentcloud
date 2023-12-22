package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumLogExportListDataSource_basic -v
func TestAccTencentCloudRumLogExportListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumLogExportListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_log_export_list.log_export_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_log_export_list.log_export_list", "result"),
				),
			},
		},
	})
}

const testAccRumLogExportListDataSource = `

data "tencentcloud_rum_log_export_list" "log_export_list" {
  project_id = 120000
}

`
