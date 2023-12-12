Use this data source to query detailed information of dbbrain slow_log_top_sqls

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_top_sqls" "test" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  sort_by = "QueryTimeMax"
  order_by = "ASC"
  product = "mysql"
}
```