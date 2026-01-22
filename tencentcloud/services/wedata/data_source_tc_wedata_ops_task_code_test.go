package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsTaskCodeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsTaskCodeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_task_code.wedata_ops_task_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_task_code.wedata_ops_task_code", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_task_code.wedata_ops_task_code", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsTaskCodeDataSource = `

data "tencentcloud_wedata_ops_task_code" "wedata_ops_task_code" {
    project_id = "1859317240494305280"
    task_id    = "20230901114849281"
}
`
