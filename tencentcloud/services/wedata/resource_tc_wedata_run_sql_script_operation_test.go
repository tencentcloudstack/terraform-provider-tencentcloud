package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataRunSqlScriptOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataRunSqlScriptOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_run_sql_script_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_run_sql_script_operation.example", "script_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_run_sql_script_operation.example", "project_id"),
				),
			},
		},
	})
}

const testAccWedataRunSqlScriptOperation = `
resource "tencentcloud_wedata_run_sql_script_operation" "example" {
  script_id  = "195a5f49-8e43-4e05-8b42-cecdfb6349f8"
  project_id = "2983848457986924544"
}
`
