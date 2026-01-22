package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataWorkflowMaxPermissionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataWorkflowMaxPermissionDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_workflow_max_permission.example"),
			),
		}},
	})
}

const testAccWedataWorkflowMaxPermissionDataSource = `
data "tencentcloud_wedata_workflow_max_permission" "example" {
  project_id  = "3108707295180644352"
  entity_id   = "53e78f97-f145-11f0-ba36-b8cef6a5af5c"
  entity_type = "folder"
}
`
