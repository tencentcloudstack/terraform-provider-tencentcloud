package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataRerunTriggerWorkflowRunAsyncResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataRerunTriggerWorkflowRunAsync,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_rerun_trigger_workflow_run_async.rerun_basic", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_rerun_trigger_workflow_run_async.rerun_basic", "project_id", "3108707295180644352"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_rerun_trigger_workflow_run_async.rerun_basic", "workflow_id", "333368d7-bc8e-4b95-9a66-7a5151063eb2"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_rerun_trigger_workflow_run_async.rerun_basic", "workflow_execution_id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_rerun_trigger_workflow_run_async.rerun_basic", "execute_type", "1"),
				),
			},
		},
	})
}

const testAccWedataRerunTriggerWorkflowRunAsync = `

data "tencentcloud_wedata_trigger_workflow_runs" "trigger_workflow_runs" {
  project_id = "3108707295180644352"
  filters {
    name   = "WorkflowId"
    values = ["333368d7-bc8e-4b95-9a66-7a5151063eb2"]
  }
  order_fields {
    name      = "CreateTime"
    direction = "DESC"
  }
}

resource "tencentcloud_wedata_rerun_trigger_workflow_run_async" "rerun_basic" {
  project_id             = "3108707295180644352"
  workflow_id            = "333368d7-bc8e-4b95-9a66-7a5151063eb2"
  workflow_execution_id  = data.tencentcloud_wedata_trigger_workflow_runs.trigger_workflow_runs.data[0].items[0].execution_id
  execute_type          = "1"
}
`
