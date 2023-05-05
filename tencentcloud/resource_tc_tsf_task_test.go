package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfTaskResource_basic -v
func TestAccTencentCloudTsfTaskResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfTask,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfTaskExists("tencentcloud_tsf_task.task"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_task.task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_content", "/test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "execute_type", "unicast"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_type", "java"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "time_out", "60000"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "group_id", defaultTsfGWGroupId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_rule.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_rule.0.rule_type", "Cron"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_rule.0.expression", "0 * 1 * * ? "),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "retry_count", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "retry_interval", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "success_operator", "GTE"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "advance_settings.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "advance_settings.0.sub_task_concurrency", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_task.task", "task_argument", "a=c"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_task.task",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
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
			if *res.TaskState == "DELETED" {
				return nil
			}
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

const testAccTsfTaskVar = `
variable "group_id" {
	default = "` + defaultTsfGWGroupId + `"
}
`

const testAccTsfTask = testAccTsfTaskVar + `

resource "tencentcloud_tsf_task" "task" {
	task_name = "terraform-test"
	task_content = "/test"
	execute_type = "unicast"
	task_type = "java"
	time_out = 60000
	group_id = var.group_id
	task_rule {
	  rule_type = "Cron"
	  expression = "0 * 1 * * ? "
	}
	retry_count = 0
	retry_interval = 0
	success_operator = "GTE"
	success_ratio = "100"
	advance_settings {
	  sub_task_concurrency = 2
	}
	task_argument = "a=c"
}

`
