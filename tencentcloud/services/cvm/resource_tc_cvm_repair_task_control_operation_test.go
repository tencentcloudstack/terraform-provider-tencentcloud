package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccTencentCloudCvmRepairTaskControlOperationResource_basic verifies the basic authorize flow.
// NOTE: This test requires a real CVM repair task in `pending authorization` state. Replace the
// `task_id` and `instance_ids` placeholders with values discoverable via the
// `tencentcloud_cvm_repair_tasks` data source before running it.
func TestAccTencentCloudCvmRepairTaskControlOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRepairTaskControlOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_repair_task_control_operation.demo", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_repair_task_control_operation.demo", "product", "CVM"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_repair_task_control_operation.demo", "operate", "AuthorizeRepair"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_repair_task_control_operation.demo", "task_id", "rep-xxxxxxxx"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_repair_task_control_operation.demo", "instance_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_repair_task_control_operation.demo", "instance_ids.0", "ins-xxxxxxxx"),
				),
			},
		},
	})
}

// TestAccTencentCloudCvmRepairTaskControlOperationResource_withOrderTime verifies the scheduled
// (order_auth_time) authorization flow.
func TestAccTencentCloudCvmRepairTaskControlOperationResource_withOrderTime(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRepairTaskControlOperationWithOrderTime,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_repair_task_control_operation.scheduled", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_repair_task_control_operation.scheduled", "order_auth_time"),
				),
			},
		},
	})
}

const testAccCvmRepairTaskControlOperation = `
resource "tencentcloud_cvm_repair_task_control_operation" "demo" {
  product      = "CVM"
  instance_ids = ["ins-xxxxxxxx"]
  task_id      = "rep-xxxxxxxx"
  operate      = "AuthorizeRepair"
}
`

const testAccCvmRepairTaskControlOperationWithOrderTime = `
resource "tencentcloud_cvm_repair_task_control_operation" "scheduled" {
  product         = "CVM"
  instance_ids    = ["ins-xxxxxxxx"]
  task_id         = "rep-xxxxxxxx"
  operate         = "AuthorizeRepair"
  order_auth_time = "2030-01-01 12:00:00"
}
`
