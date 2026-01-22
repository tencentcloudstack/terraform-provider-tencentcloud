package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_task.wedata_ops_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task.wedata_ops_task", "action", "PAUSE"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task.wedata_ops_task", "status", "O"),
				),
			},
			{
				Config: testAccWedataOpsTaskUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_task.wedata_ops_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task.wedata_ops_task", "action", "START"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_ops_task.wedata_ops_task", "status", "Y"),
				),
			},
		},
	})
}

const testAccWedataOpsTask = `

resource "tencentcloud_wedata_ops_task" "wedata_ops_task" {
  project_id = "1859317240494305280"
  task_id    = "20251013154418424"
  action     = "PAUSE"
}
`

const testAccWedataOpsTaskUp = `

resource "tencentcloud_wedata_ops_task" "wedata_ops_task" {
  project_id = "1859317240494305280"
  task_id    = "20251013154418424"
  action     = "START"
}
`
