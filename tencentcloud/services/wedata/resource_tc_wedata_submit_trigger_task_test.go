package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataSubmitTriggerTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataSubmitTriggerTask,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_wedata_submit_trigger_task.wedata_submit_trigger_task", "id"),
				resource.TestCheckResourceAttr("tencentcloud_wedata_submit_trigger_task.wedata_submit_trigger_task", "project_id", "3108707295180644352"),
				resource.TestCheckResourceAttr("tencentcloud_wedata_submit_trigger_task.wedata_submit_trigger_task", "task_id", "20260115150233414"),
				resource.TestCheckResourceAttr("tencentcloud_wedata_submit_trigger_task.wedata_submit_trigger_task", "version_remark", "v1"),
			),
		}},
	})
}

const testAccWedataSubmitTriggerTask = `

resource "tencentcloud_wedata_submit_trigger_task" "wedata_submit_trigger_task" {
	project_id = "3108707295180644352"
	task_id = "20260115150233414"
	version_remark = "v1"
}
`
