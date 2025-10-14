package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsWorkflowDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsWorkflowDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_workflow.wedata_ops_workflow"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_workflow.wedata_ops_workflow", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_workflow.wedata_ops_workflow", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsWorkflowDataSource = `

data "tencentcloud_wedata_ops_workflow" "wedata_ops_workflow" {
    project_id = "2905622749543821312"
    workflow_id = "f328ab83-62e1-4b0a-9a18-a79b42722792"
}
`
