package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testExecutionTasksObjectName = "data.tencentcloud_tcr_tag_retention_execution_tasks.tasks"

func TestAccTencentCloudTcrTagRetentionExecutionTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionExecutionTasksDataSource,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testExecutionTasksObjectName),
					resource.TestCheckResourceAttrSet(testExecutionTasksObjectName, "id"),
					resource.TestCheckResourceAttrSet(testExecutionTasksObjectName, "registry_id"),
					resource.TestCheckResourceAttr(testExecutionTasksObjectName, "retention_id", "1"),
					resource.TestCheckResourceAttr(testExecutionTasksObjectName, "execution_id", "1"),
					resource.TestCheckResourceAttrSet(testExecutionTasksObjectName, "retention_task_list.#"),
					resource.TestCheckResourceAttrSet(testExecutionTasksObjectName, "retention_task_list.0.task_id"),
					resource.TestCheckResourceAttrSet(testExecutionTasksObjectName, "retention_task_list.0.execution_id"),
				),
			},
		},
	})
}

const testAccTcrTagRetentionExecutionTasksDataSource = TCRDataSource + `

data "tencentcloud_tcr_tag_retention_execution_tasks" "tasks" {
  registry_id = local.tcr_id
  retention_id = 1
  execution_id = 1
}

`
