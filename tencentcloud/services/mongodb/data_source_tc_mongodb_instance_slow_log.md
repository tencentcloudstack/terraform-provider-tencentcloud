Use this data source to query detailed information of mongodb instance_slow_log

Example Usage

```hcl
data "tencentcloud_mongodb_instance_slow_log" "instance_slow_log" {
  instance_id = "cmgo-9d0p6umb"
  start_time = "2019-06-01 10:00:00"
  end_time = "2019-06-02 12:00:00"
  slow_m_s = 100
  format = "json"
}
```