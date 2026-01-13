package wedata_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudWedataTriggerTaskRunDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerTaskRunDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_task_run.trigger_task_run"),
				),
			},
		},
	})
}

const testAccWedataTriggerTaskRunDataSource = `
data "tencentcloud_wedata_trigger_task_run" "trigger_task_run" {
  project_id = "1840731342293818368"
  task_execution_id = "20241205_2057_1840731342293818368_20241205205700_workflow_20241205205700_20241205205700_1840731342293818368_20241205205700_task_20241205205700_20241205205700_1840731342293818368_20241205205700"
}
`
