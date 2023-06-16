package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrReplicationInstanceCreateTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Config: fmt.Sprintf(testAccTcrReplicationInstance_create_tasks_and_sync_status_DataSource, "create-tasks-and-sync-status", time.Now().Nanosecond()),
				Config: testAccTcrReplicationInstance_create_tasks_and_sync_status_DataSource,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "replication_registry_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "replication_region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "status"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.0.task_name"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.0.task_uuid"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.0.task_status"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.0.task_message"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.0.created_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_create_tasks.create_tasks", "task_detail.0.finished_time"),

					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_replication_instance_sync_status.sync_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "registry_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_registry_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_region_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "show_replication_log", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.0.resource_type"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.0.source"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.0.destination"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.0.status"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.0.start_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_replication_instance_sync_status.sync_status", "replication_log.0.end_time"),
				),
			},
		},
	})
}

// const testAccTcrReplicationInstance_create_tasks_and_sync_status_DataSource = testAccTcrManageReplicationOperation + `
const testAccTcrReplicationInstance_create_tasks_and_sync_status_DataSource = `
// locals {
//   src_registry_id = local.tcr_id
//   dst_registry_id = tencentcloud_tcr_manage_replication_operation.my_replica.destination_registry_id
//   dst_region_id   = tencentcloud_tcr_manage_replication_operation.my_replica.destination_region_id
// }
locals {
  src_registry_id = "tcr-800uf6j9"
  dst_registry_id = "tcr-800uf6j9-1-cgeafk"
  dst_region_id   = 1
}

data "tencentcloud_tcr_replication_instance_create_tasks" "create_tasks" {
  replication_registry_id = local.dst_registry_id
  replication_region_id   = local.dst_region_id
}

data "tencentcloud_tcr_replication_instance_sync_status" "sync_status" {
  registry_id             = local.src_registry_id
  replication_registry_id = local.dst_registry_id
  replication_region_id   = local.dst_region_id
  show_replication_log    = false
}

`
