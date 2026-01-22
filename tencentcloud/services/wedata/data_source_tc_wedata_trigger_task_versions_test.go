package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerTaskVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerTaskVersionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "project_id", "3108707295180644352"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "task_id", "20260109011325782"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.#"),
					// Check version items if they exist
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.version_num"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.create_user_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.version_id"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.version_remark"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.approve_status"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_task_versions.wedata_trigger_task_versions", "data.0.items.0.status"),
				),
			},
		},
	})
}

const testAccWedataTriggerTaskVersionsDataSource = `

data "tencentcloud_wedata_trigger_task_versions" "wedata_trigger_task_versions" {
  project_id = "3108707295180644352"
  task_id    = "20260109011325782"
}
`
