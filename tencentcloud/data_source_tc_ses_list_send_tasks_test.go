package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesListSendTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesListSendTasksDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_list_send_tasks.list_send_tasks")),
			},
		},
	})
}

const testAccSesListSendTasksDataSource = `

data "tencentcloud_ses_list_send_tasks" "list_send_tasks" {
  status = 1
  receiver_id = 124
  task_type = 2
  }

`
