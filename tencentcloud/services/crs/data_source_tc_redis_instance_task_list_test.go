package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisInstanceTaskListDataSource_basic -v
func TestAccTencentCloudRedisInstanceTaskListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceTaskListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_task_list.instance_task_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.progress"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.result"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.task_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_instance_task_list.instance_task_list", "tasks.0.task_type"),
				),
			},
		},
	})
}

const testAccRedisInstanceTaskListDataSource = `

data "tencentcloud_redis_instance_task_list" "instance_task_list" {
	instance_id = "crs-jf4ico4v"
	instance_name = "Keep-terraform"
	project_ids = [0]
	task_types = ["034"]
	begin_time = "2023-04-09 23:03:31"
	end_time = "2023-04-09 23:03:51"
	task_status = [2]
	result = [2]
	# operate_uin = [""]
}

`
