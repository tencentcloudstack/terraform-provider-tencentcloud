Use this data source to query detailed information of dbbrain diag_db_instances

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_db_instances" "diag_db_instances" {
	is_supported   = true
	product        = "mysql"
	instance_names = ["keep_preset_mysql"]
}
```