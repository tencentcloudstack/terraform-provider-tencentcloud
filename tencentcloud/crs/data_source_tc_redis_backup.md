Use this data source to query detailed information of redis backup

Example Usage

```hcl
data "tencentcloud_redis_backup" "backup" {
  instance_id = "crs-c1nl9rpv"
  begin_time = "2023-04-07 03:57:30"
  end_time = "2023-04-07 03:57:56"
  status = [2]
  instance_name = "Keep-terraform"
}
```