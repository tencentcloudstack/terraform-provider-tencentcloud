package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsUpstreamTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsUpstreamTasksDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_upstream_tasks.wedata_ops_upstream_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_upstream_tasks.wedata_ops_upstream_tasks", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_upstream_tasks.wedata_ops_upstream_tasks", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsUpstreamTasksDataSource = `

data "tencentcloud_wedata_ops_upstream_tasks" "wedata_ops_upstream_tasks" {
    project_id = "1859317240494305280"
    task_id    = "20250820150144998"
}
`
