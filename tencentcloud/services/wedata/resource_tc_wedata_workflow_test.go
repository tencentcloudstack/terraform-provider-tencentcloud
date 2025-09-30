package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataWorkflowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataWorkflow,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_workflow.wedata_workflow", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_workflow.wedata_workflow",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataWorkflow = `
resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id = 2905622749543821312
  workflow_name = "test"
  parent_folder_path = "/tfmika"
  workflow_type = "cycle"
}
`
