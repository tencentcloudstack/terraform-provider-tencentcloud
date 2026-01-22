package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsWorkflowsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsWorkflowsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_workflows.wedata_ops_workflows"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_workflows.wedata_ops_workflows", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_workflows.wedata_ops_workflows", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataOpsWorkflowsDataSource = `

data "tencentcloud_wedata_ops_workflows" "wedata_ops_workflows" {
    project_id = "2905622749543821312"
    folder_id = "720ecbfb-7e5a-11f0-ba36-b8cef6a5af5c"
    status = "ALL_RUNNING"
    owner_uin = "100044349576"
    workflow_type = "Cycle"
    sort_type = "ASC"
}
`
