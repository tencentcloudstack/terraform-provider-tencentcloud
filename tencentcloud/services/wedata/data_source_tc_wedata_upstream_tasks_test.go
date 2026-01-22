package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataUpstreamTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataUpstreamTasksDataSource,
			Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_upstream_tasks.wedata_upstream_tasks")),
		}},
	})
}

const testAccWedataUpstreamTasksDataSource = `

data "tencentcloud_wedata_upstream_tasks" "wedata_upstream_tasks" {
	project_id = "2905622749543821312"
	task_id = "20251015164958429"
}
`
