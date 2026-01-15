package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsTriggerWorkflowsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataOpsTriggerWorkflowsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows", "data.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows", "data.0.total_count"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows", "data.0.total_page_number"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows", "data.0.page_number"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows", "data.0.page_size"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_trigger_workflows.wedata_ops_trigger_workflows", "data.0.items.#"),
			),
		}},
	})
}

const testAccWedataOpsTriggerWorkflowsDataSource = `
data "tencentcloud_wedata_ops_trigger_workflows" "wedata_ops_trigger_workflows" {
  project_id = "3108707295180644352"
}
`
