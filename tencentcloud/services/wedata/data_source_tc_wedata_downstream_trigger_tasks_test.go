package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataDownstreamTriggerTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDownstreamTriggerTasksDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "project_id", "3108707295180644352"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "task_id", "20260109165716558"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.#"),
					// Check task items if they exist
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.task_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.task_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.workflow_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.workflow_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.task_type_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.task_type_desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_downstream_trigger_tasks.wedata_downstream_trigger_tasks", "data.0.items.0.owner_uin"),
				),
			},
		},
	})
}

const testAccWedataDownstreamTriggerTasksDataSource = `

data "tencentcloud_wedata_downstream_trigger_tasks" "wedata_downstream_trigger_tasks" {
  project_id = "3108707295180644352"
  task_id    = "20260109165716558"
}
`
