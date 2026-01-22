package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsTasksDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_tasks.wedata_ops_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_tasks.wedata_ops_tasks", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_tasks.wedata_ops_tasks", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsTasksDataSource = `

data "tencentcloud_wedata_ops_tasks" "wedata_ops_tasks" {
  project_id        = "1859317240494305280"
  task_type_id      = 34
  workflow_id       = "d7184172-4879-11ee-ba36-b8cef6a5af5c"
  workflow_name     = "test1"
  folder_id         = "cee5780a-4879-11ee-ba36-b8cef6a5af5c"
  executor_group_id = "20230830105723839685"
  cycle_type        = "MINUTE_CYCLE"
  status            = "F"
}
`
