package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataIntegrationOfflineTaskResource_basic -v
func TestAccTencentCloudNeedFixWedataIntegrationOfflineTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegrationOfflineTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_offline_task.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "end_time", "2099-12-31 00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "notes", "terraform example demo."),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "start_time", "2023-12-31 00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "task_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "task_action", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "task_mode", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_offline_task.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataIntegrationOfflineTaskUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_offline_task.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "end_time", "2199-12-31 00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "notes", "terraform example demo."),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "start_time", "2024-12-31 00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "task_name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "task_action", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_offline_task.example", "task_mode", "1"),
				),
			},
		},
	})
}

const testAccWedataIntegrationOfflineTask = `
resource "tencentcloud_wedata_integration_offline_task" "example" {
  project_id  = "1612982498218618880"
  cycle_step  = 1
  delay_time  = 0
  end_time    = "2099-12-31 00:00:00"
  notes       = "terraform example demo."
  start_time  = "2023-12-31 00:00:00"
  task_name   = "tf_example"
  task_action = "2"
  task_mode   = "1"

  task_info {
    executor_id = "20230313175748567418"
    config {
      name  = "Args"
      value = "args"
    }
    config {
      name  = "dirtyDataThreshold"
      value = "0"
    }
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "syncRateLimitUnit"
      value = "0"
    }
    ext_config {
      name  = "TaskAlarmRegularList"
      value = "73"
    }
    incharge = "tiyan1"
    offline_task_add_entity {
      cycle_type         = 3
      crontab_expression = "0 0 1 * * ?"
      retry_wait         = 5
      retriable          = 1
      try_limit          = 5
      self_depend        = 1
    }
  }
}
`

const testAccWedataIntegrationOfflineTaskUpdate = `
resource "tencentcloud_wedata_integration_offline_task" "example" {
  project_id  = "1612982498218618880"
  cycle_step  = 1
  delay_time  = 0
  end_time    = "2199-12-31 00:00:00"
  notes       = "terraform example demo."
  start_time  = "2024-12-31 00:00:00"
  task_name   = "tf_example_update"
  task_action = "2"
  task_mode   = "1"

  task_info {
    executor_id = "20230313175748567418"
    config {
      name  = "Args"
      value = "args"
    }
    config {
      name  = "dirtyDataThreshold"
      value = "0"
    }
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "syncRateLimitUnit"
      value = "0"
    }
    ext_config {
      name  = "TaskAlarmRegularList"
      value = "73"
    }
    incharge = "tiyan1"
    offline_task_add_entity {
      cycle_type         = 3
      crontab_expression = "0 0 1 * * ?"
      retry_wait         = 5
      retriable          = 1
      try_limit          = 5
      self_depend        = 1
    }
  }
}
`
