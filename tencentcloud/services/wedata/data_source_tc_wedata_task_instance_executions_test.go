package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTaskInstanceExecutionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTaskInstanceExecutionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_task_instance_executions.wedata_task_instance_executions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_task_instance_executions.wedata_task_instance_executions", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_task_instance_executions.wedata_task_instance_executions", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataTaskInstanceExecutionsDataSource = `

data "tencentcloud_wedata_task_instance_executions" "wedata_task_instance_executions" {
  project_id   = "1859317240494305280"
  instance_key = "20250731151633120_2025-10-13 17:00:00"
}
`
