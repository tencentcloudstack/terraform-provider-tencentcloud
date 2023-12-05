Provides a resource to create a dbbrain security_audit_log_export_task

Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "sec_audit_group_id"
  start_time = "2020-12-28 00:00:00"
  end_time = "2020-12-28 01:00:00"
  product = "mysql"
  danger_levels = [0,1,2]
}

```