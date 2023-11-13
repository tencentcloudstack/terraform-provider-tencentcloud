package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrTagRetentionExecutionTaskDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionExecutionTaskDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_tag_retention_execution_task.tag_retention_execution_task")),
			},
		},
	})
}

const testAccTcrTagRetentionExecutionTaskDataSource = `

data "tencentcloud_tcr_tag_retention_execution_task" "tag_retention_execution_task" {
  registry_id = "tcr-xxx"
  retention_id = 1
  execution_id = 1
  }

`
