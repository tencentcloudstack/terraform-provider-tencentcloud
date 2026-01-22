package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerWorkflowRunsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataTriggerWorkflowRunsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_workflow_runs.wedata_trigger_workflow_runs"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflow_runs.wedata_trigger_workflow_runs", "data.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflow_runs.wedata_trigger_workflow_runs", "data.0.items.#"),
			),
		}},
	})
}

const testAccWedataTriggerWorkflowRunsDataSource = `
data "tencentcloud_wedata_trigger_workflow_runs" "wedata_trigger_workflow_runs" {
  project_id = "3108707295180644352"
}
`
