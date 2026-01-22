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
		Steps: []resource.TestStep{
			{
				Config: testAccWedataWorkflow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_workflow.wedata_workflow", "id")),
			},
			{
				Config: testAccWedataWorkflowUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_wedata_workflow.wedata_workflow", "workflow_name", "test1"),
				),
			},
			{
				ResourceName:            "tencentcloud_wedata_workflow.wedata_workflow",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workflow_scheduler_configuration.0.start_time"},
			}},
	})
}

const testAccWedataWorkflow = `
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
  folder_name        = "tftest"
}

resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id = 2905622749543821312
  workflow_name = "test"
  parent_folder_path = "${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.parent_folder_path}${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.folder_name}"
  workflow_type = "cycle"
}
`

const testAccWedataWorkflowUpdate = `
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
  folder_name        = "tftest"
}

resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id = 2905622749543821312
  workflow_name = "test1"
  parent_folder_path = "${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.parent_folder_path}${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.folder_name}"
  workflow_type = "cycle"
}
`
