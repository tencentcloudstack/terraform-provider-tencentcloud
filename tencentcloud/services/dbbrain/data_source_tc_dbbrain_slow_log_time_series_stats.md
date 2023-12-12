Use this data source to query detailed information of dbbrain slow_log_time_series_stats

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_time_series_stats" "test" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
}
```