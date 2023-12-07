Use this data source to query detailed information of dbbrain security_audit_log_download_urls

Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
	sec_audit_group_id = "%s"
	start_time = "%s"
	end_time = "%s"
	product = "mysql"
	danger_levels = [0,1,2]
}

data "tencentcloud_dbbrain_security_audit_log_download_urls" "test" {
	sec_audit_group_id = "%s"
	async_request_id = tencentcloud_dbbrain_security_audit_log_export_task.task.async_request_id
	product = "mysql"
}
```