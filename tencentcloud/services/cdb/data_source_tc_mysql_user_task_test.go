package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlUserTaskDataSource_basic -v
func TestAccTencentCloudMysqlUserTaskDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlUserTaskDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_user_task.user_task"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.async_request_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.instance_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.job_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.message"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.progress"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.task_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_user_task.user_task", "items.0.task_type"),
				),
			},
		},
	})
}

const testAccMysqlUserTaskDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`
const testAccMysqlUserTaskDataSource = testAccMysqlUserTaskDataSourceVar + `

data "tencentcloud_mysql_user_task" "user_task" {
	instance_id = "cdb-fitq5t9h"
}

`
