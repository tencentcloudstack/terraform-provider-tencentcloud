package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDbbrainSecurityAuditLogExportTasksDataSource(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().Add(-2 * time.Hour).In(loc).Format("2006-01-02T15:04:05+08:00")
	endTime := time.Now().Add(2 * time.Hour).In(loc).Format("2006-01-02T15:04:05+08:00")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbbrainSecurityAuditLogExportTasks(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_security_audit_log_export_tasks.tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_security_audit_log_export_tasks.tasks", "list.#"),
				),
			},
		},
	})
}

func testAccDataSourceDbbrainSecurityAuditLogExportTasks(st, et string) string {
	return fmt.Sprintf(`

resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
  danger_levels = [0,1,2]
}

data "tencentcloud_dbbrain_security_audit_log_export_tasks" "tasks" {
	sec_audit_group_id = "%s"
	product = "mysql"
	async_request_ids = [tencentcloud_dbbrain_security_audit_log_export_task.task.async_request_id]
}

`, defaultDbBrainsagId, st, et, defaultDbBrainsagId)
}
