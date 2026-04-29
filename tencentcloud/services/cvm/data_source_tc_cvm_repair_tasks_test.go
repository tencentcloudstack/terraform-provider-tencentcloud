package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmRepairTasksDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRepairTasksDataSource_Basic,
				Check: resource.ComposeTestCheckFunc(
					acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_repair_tasks.tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_repair_tasks.tasks", "repair_task_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_repair_tasks.tasks", "total_count"),
				),
			},
		},
	})
}

func TestAccTencentCloudCvmRepairTasksDataSource_Filter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRepairTasksDataSource_Filter,
				Check: resource.ComposeTestCheckFunc(
					acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_repair_tasks.filtered"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_repair_tasks.filtered", "repair_task_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_repair_tasks.filtered", "total_count"),
				),
			},
		},
	})
}

const testAccCvmRepairTasksDataSource_Basic = `
data "tencentcloud_cvm_repair_tasks" "tasks" {
}
`

const testAccCvmRepairTasksDataSource_Filter = `
data "tencentcloud_cvm_repair_tasks" "filtered" {
  product = "CVM"
  task_status = [1, 2, 4]
}
`
