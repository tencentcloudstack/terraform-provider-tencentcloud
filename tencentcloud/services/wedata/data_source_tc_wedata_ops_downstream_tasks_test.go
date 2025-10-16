package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsDownstreamTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsDownstreamTasksDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_downstream_tasks.wedata_ops_downstream_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_downstream_tasks.wedata_ops_downstream_tasks", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_downstream_tasks.wedata_ops_downstream_tasks", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsDownstreamTasksDataSource = `

data "tencentcloud_wedata_ops_downstream_tasks" "wedata_ops_downstream_tasks" {
    project_id = "1859317240494305280"
    task_id    = "20250820150144998"
}
`
