Provides a resource to create a dbbrain db_diag_report_task

Example Usage

```hcl
resource "tencentcloud_dbbrain_db_diag_report_task" "db_diag_report_task" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  send_mail_flag = 0
  product = "mysql"
}
```