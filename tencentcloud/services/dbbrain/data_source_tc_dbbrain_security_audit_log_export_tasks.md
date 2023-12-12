Use this data source to query detailed information of dbbrain securityAuditLogExportTasks

Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "sec_audit_group_id"
  start_time = "start_time"
  end_time = "end_time"
  product = "mysql"
  danger_levels = [0,1,2]
}

data "tencentcloud_dbbrain_security_audit_log_export_tasks" "tasks" {
	sec_audit_group_id = "sec_audit_group_id"
	product = "mysql"
	async_request_ids = [tencentcloud_dbbrain_security_audit_log_export_task.task.async_request_id]
}
```