package tencentcloud

import (
	"fmt"
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
				Config: fmt.Sprintf(testAccTcrTagRetentionExecutionTasksDataSource, defaultTCRInstanceId),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testExecutionTasksObjectName),
					resource.TestCheckResourceAttrSet(testExecutionTasksObjectName, "id"),
					resource.TestCheckResourceAttr(testExecutionTasksObjectName, "registry_id", defaultTCRInstanceId),
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

const testAccTcrTagRetentionExecutionTasksDataSource = `

data "tencentcloud_tcr_tag_retention_execution_tasks" "tasks" {
  registry_id = "%s"
  retention_id = 1
  execution_id = 1
}

`
