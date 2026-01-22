package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataStopSqlScriptRunOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataStopSqlScriptRunOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_stop_sql_script_run_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_stop_sql_script_run_operation.example", "job_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_stop_sql_script_run_operation.example", "project_id"),
				),
			},
		},
	})
}

const testAccWedataStopSqlScriptRunOperation = `
resource "tencentcloud_wedata_stop_sql_script_run_operation" "example" {
  job_id     = "ac13aceb-7a30-4414-91c0-6504f177462f"
  project_id = "2983848457986924544"
}
`
