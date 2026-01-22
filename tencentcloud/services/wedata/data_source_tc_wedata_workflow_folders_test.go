package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataWorkflowFoldersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataWorkflowFoldersDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_workflow_folders.wedata_workflow_folders"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_workflow_folders.wedata_workflow_folders", "data.#"),
			),
		}},
	})
}

const testAccWedataWorkflowFoldersDataSource = `
data "tencentcloud_wedata_workflow_folders" "wedata_workflow_folders" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
}
`
