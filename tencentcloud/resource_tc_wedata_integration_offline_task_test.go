package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataIntegration_offline_taskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegration_offline_task,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_offline_task.integration_offline_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_offline_task.integration_offline_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataIntegration_offline_task = `

resource "tencentcloud_wedata_integration_offline_task" "integration_offline_task" {
  project_id = "1455251608631480391"
  cycle_step = 1
  delay_time = 0
  end_time = "2099-12-31 00:00:00"
  notes = "Task for test"
  start_time = "2023-12-31 00:00:00"
  task_name = "TaskTest_10"
  type_id = 27
  task_action = "0,3,4"
  task_mode = "1"
}

`
