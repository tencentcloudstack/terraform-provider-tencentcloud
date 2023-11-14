package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisInstanceTaskListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceTaskListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_task_list.instance_task_list")),
			},
		},
	})
}

const testAccRedisInstanceTaskListDataSource = `

data "tencentcloud_redis_instance_task_list" "instance_task_list" {
  instance_id = "crs-c1nl9rpv"
  instance_name = &lt;nil&gt;
  project_ids = 
  task_types = 
  begin_time = "2021-12-30 00:00:00"
  end_time = "2021-12-30 00:00:00"
  task_status = &lt;nil&gt;
  result = &lt;nil&gt;
  operate_uin = 
  }

`
