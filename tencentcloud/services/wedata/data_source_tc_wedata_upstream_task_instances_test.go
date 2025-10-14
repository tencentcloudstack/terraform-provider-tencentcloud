package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataUpstreamTaskInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataUpstreamTaskInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_upstream_task_instances.wedata_upstream_task_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_upstream_task_instances.wedata_upstream_task_instances", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_upstream_task_instances.wedata_upstream_task_instances", "data.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_upstream_task_instances.wedata_upstream_task_instances", "data.0.items.#", "2"),
				),
			},
		},
	})
}

const testAccWedataUpstreamTaskInstancesDataSource = `

data "tencentcloud_wedata_upstream_task_instances" "wedata_upstream_task_instances" {
  project_id   = "1859317240494305280"
  instance_key = "20250820150153594_2025-10-13 17:00:00"
}
`
