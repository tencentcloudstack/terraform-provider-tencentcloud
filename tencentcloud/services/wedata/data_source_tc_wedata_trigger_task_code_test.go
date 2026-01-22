package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerTaskCodeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerTaskCodeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code", "project_id", "3108707295180644352"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code", "task_id", "20260109011325782"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code", "data.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code", "data.0.code_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_code.wedata_trigger_task_code", "data.0.code_file_size"),
				),
			},
		},
	})
}

const testAccWedataTriggerTaskCodeDataSource = `

data "tencentcloud_wedata_trigger_task_code" "wedata_trigger_task_code" {
  project_id = "3108707295180644352"
  task_id    = "20260109011325782"
}
`
