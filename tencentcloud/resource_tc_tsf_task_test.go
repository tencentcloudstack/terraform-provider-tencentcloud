package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTsfTaskResource_basic -v
func TestAccTencentCloudNeedFixTsfTaskResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfTask,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfTaskExists("tencentcloud_tsf_task.task"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_task.task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_name", ""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_content", ""),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_task.task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfTaskDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_task" {
			continue
		}

		res, err := service.DescribeTsfTaskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf Task %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfTaskExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfTaskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf Task %s is not found", rs.Primary.ID)
		}

		return nil
	}
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
