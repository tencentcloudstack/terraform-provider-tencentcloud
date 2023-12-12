Provides a resource to create a dbbrain modify_diag_db_instance_conf

Example Usage

```hcl
resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "on" {
  instance_confs {
	daily_inspection = "Yes"
	overview_display = "Yes"
  }
  product = "mysql"
  instance_ids = ["%s"]
}
```

```hcl
resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "off" {
  instance_confs {
	daily_inspection = "No"
	overview_display = "No"
  }
  product = "mysql"
  instance_ids = ["%s"]
}
```