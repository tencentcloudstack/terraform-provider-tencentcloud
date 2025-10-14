package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTaskInstanceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTaskInstanceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_task_instance.wedata_task_instance"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_task_instance.wedata_task_instance", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_task_instance.wedata_task_instance", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataTaskInstanceDataSource = `

data "tencentcloud_wedata_task_instance" "wedata_task_instance" {
  project_id = "1859317240494305280"
  instance_key = "20250324192240178_2025-10-13 11:50:00"
}
`
