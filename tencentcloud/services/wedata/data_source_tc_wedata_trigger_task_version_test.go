package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerTaskVersionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerTaskVersionDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "project_id", "3108707295180644352"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "task_id", "20260109011325782"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.#", "1"),
					// Check version details
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.version_num"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.create_user_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.version_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.version_remark"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.approve_status"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.approve_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.approve_user_uin"),
					// Check task details
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.task.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.task.0.trigger_task_base_attribute.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.task.0.trigger_task_base_attribute.0.task_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.task.0.trigger_task_base_attribute.0.task_type_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_version.wedata_trigger_task_version", "data.0.task.0.trigger_task_base_attribute.0.workflow_id"),
				),
			},
		},
	})
}

const testAccWedataTriggerTaskVersionDataSource = `

data "tencentcloud_wedata_trigger_task_version" "wedata_trigger_task_version" {
  project_id = "3108707295180644352"
  task_id    = "20260109011325782"
}
`
