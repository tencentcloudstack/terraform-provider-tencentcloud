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
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_offline_task.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataIntegrationOfflineTask = `
resource "tencentcloud_wedata_integration_offline_task" "example" {
  project_id  = "1455251608631480391"
  cycle_step  = 1
  delay_time  = 0
  end_time    = "2099-12-31 00:00:00"
  notes       = "terraform example demo."
  start_time  = "2023-12-31 00:00:00"
  task_name   = "tf_example"
  type_id     = 27
  task_action = "0,3,4"
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
    incharge = "demo_user"
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
