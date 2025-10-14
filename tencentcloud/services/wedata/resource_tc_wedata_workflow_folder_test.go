package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataWorkflowFolderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataWorkflowFolder,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_workflow_folder.wedata_workflow_folder", "id")),
			},
			{
				Config: testAccWedataWorkflowFolderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_wedata_workflow_folder.wedata_workflow_folder", "folder_name", "tftest1"),
				),
			},
		},
	})
}

const testAccWedataWorkflowFolder = `
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
  folder_name        = "tftest"
}
`

const testAccWedataWorkflowFolderUpdate = `
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
  folder_name        = "tftest1"
}
`
