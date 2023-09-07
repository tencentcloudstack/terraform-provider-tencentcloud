package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudSesSendTasksDataSource_basic -v
func TestAccTencentCloudSesSendTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesSendTasksDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_send_tasks.send_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.cache_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.from_email_address"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.receivers_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.request_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.send_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.subject"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.task_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.task_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.task_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.template.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.template.0.template_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_tasks.send_tasks", "data.0.update_time"),
				),
			},
		},
	})
}

const testAccSesSendTasksDataSource = `

data "tencentcloud_ses_send_tasks" "send_tasks" {
  status = 10
  receiver_id = 1063742
  task_type = 1
}

`
