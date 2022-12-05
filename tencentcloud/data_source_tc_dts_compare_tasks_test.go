package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsCompareTasksDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDtsCompareTasks, defaultDTSJobId, defaultDTSJobId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dts_compare_tasks.compare_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_compare_tasks.compare_tasks", "list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDtsCompareTasks = `
resource "tencentcloud_dts_compare_task" "task" {
	job_id = "%s"
	task_name = "tf_test_compare_task"
	objects {
	  object_mode = "all"
	}
  }

data "tencentcloud_dts_compare_tasks" "compare_tasks" {
  job_id = "%s"
  }

`
