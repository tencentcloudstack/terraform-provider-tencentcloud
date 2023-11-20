package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataIntegrationRealtimeTaskResource_basic -v
func TestAccTencentCloudNeedFixWedataIntegrationRealtimeTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegrationRealtimeTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "task_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "task_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "description", "description."),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "sync_type"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_realtime_task.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataIntegrationRealtimeTaskUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "task_name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "task_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_realtime_task.example", "description", "description update."),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "sync_type"),
				),
			},
		},
	})
}

const testAccWedataIntegrationRealtimeTask = `
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1612982498218618880"
  task_name   = "tf_example"
  task_mode   = "1"
  description = "description."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230313175748567418"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = "35"
    }
  }
}
`

const testAccWedataIntegrationRealtimeTaskUpdate = `
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1612982498218618880"
  task_name   = "tf_example_update"
  task_mode   = "1"
  description = "description update."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230313175748567418"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = "35"
    }
  }
}
`
