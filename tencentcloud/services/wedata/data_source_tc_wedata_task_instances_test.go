package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTaskInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTaskInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_task_instances.wedata_task_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_task_instances.wedata_task_instances", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_task_instances.wedata_task_instances", "data.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_task_instances.wedata_task_instances", "data.0.items.#", "10"),
				),
			},
		},
	})
}

const testAccWedataTaskInstancesDataSource = `

data "tencentcloud_wedata_task_instances" "wedata_task_instances" {
  project_id = "1859317240494305280"
}
`
