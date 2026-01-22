package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataSqlScriptRunsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataSqlScriptRunsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_sql_script_runs.example"),
				resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script_runs.example", "id"),
				resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script_runs.example", "project_id"),
				resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script_runs.example", "script_id"),
			),
		}},
	})
}

const testAccWedataSqlScriptRunsDataSource = `
data "tencentcloud_wedata_sql_script_runs" "example" {
  project_id = "1460947878944567296"
  script_id  = "971c1520-836f-41be-b13f-7a6c637317c8"
}
`
