package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbUserTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbUserTasksDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_user_tasks.user_tasks")),
			},
		},
	})
}

const testAccDcdbUserTasksDataSource = `

data "tencentcloud_dcdb_user_tasks" "user_tasks" {
  statuses = 
  instance_ids = 
  flow_types = 
  start_time = ""
  end_time = ""
  u_task_ids = 
  }

`
