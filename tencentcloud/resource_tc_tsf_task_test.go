package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_task.task", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_task.task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfTask = `

resource "tencentcloud_tsf_task" "task" {
    task_name = ""
  task_content = ""
  execute_type = ""
  task_type = ""
  time_out = 
  group_id = ""
  task_rule {
		rule_type = ""
		expression = ""
		repeat_interval = 

  }
  retry_count = 
  retry_interval = 
  shard_count = 
  shard_arguments {
		shard_key = 
		shard_value = ""

  }
  success_operator = ""
  success_ratio = ""
  advance_settings {
		sub_task_concurrency = 

  }
  task_argument = ""
          program_id_list = 
}

`
