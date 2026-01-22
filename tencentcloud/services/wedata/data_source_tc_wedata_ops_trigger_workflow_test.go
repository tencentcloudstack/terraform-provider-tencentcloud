package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsTriggerWorkflowDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataOpsTriggerWorkflowDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_trigger_workflow.wedata_ops_trigger_workflow"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflow.wedata_ops_trigger_workflow", "data.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflow.wedata_ops_trigger_workflow", "data.0.trigger_tasks.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflow.wedata_ops_trigger_workflow", "data.0.trigger_task_links.#"),
			),
		}},
	})
}

const testAccWedataOpsTriggerWorkflowDataSource = `
data "tencentcloud_wedata_ops_trigger_workflow" "wedata_ops_trigger_workflow" {
  project_id = "3108707295180644352"
  workflow_id = "b41e8d13-905a-4006-9d05-1fe180338f59"
}
`
