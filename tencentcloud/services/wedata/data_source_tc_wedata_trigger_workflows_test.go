package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerWorkflowsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerWorkflowsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "project_id", "3108707295180644352"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.#"),
					// Check workflow items if they exist
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.workflow_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.workflow_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.owner_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.update_user_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_trigger_workflows.wedata_trigger_workflows", "data.0.items.0.create_user_uin"),
				),
			},
		},
	})
}

const testAccWedataTriggerWorkflowsDataSource = `

data "tencentcloud_wedata_trigger_workflows" "wedata_trigger_workflows" {
  project_id = "3108707295180644352"
}
`
