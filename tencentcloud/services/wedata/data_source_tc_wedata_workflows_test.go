package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataWorkflowsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataWorkflowsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_workflows.wedata_workflows"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_workflows.wedata_workflows", "data.#"),
			),
		}},
	})
}

const testAccWedataWorkflowsDataSource = `
data "tencentcloud_wedata_workflows" "wedata_workflows" {
  project_id = 2905622749543821312
  keyword    = "test_workflow"
}
`
