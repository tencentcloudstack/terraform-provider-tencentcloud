package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerWorkflowRunDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerWorkflowRunDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_workflow_run.trigger_workflow_run"),
				),
			},
		},
	})
}

const testAccWedataTriggerWorkflowRunDataSource = `
data "tencentcloud_wedata_trigger_workflow_run" "trigger_workflow_run" {
  project_id = "3108707295180644352"
  workflow_execution_id = "4087f932-9e04-41f9-879a-24fbab7f8789"
}
`
