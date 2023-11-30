Use this data source to query detailed information of dbbrain slow_log_user_sql_advice

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_user_sql_advice" "test" {
  instance_id = "%s"
  sql_text = "%s"
  product = "mysql"
}
```