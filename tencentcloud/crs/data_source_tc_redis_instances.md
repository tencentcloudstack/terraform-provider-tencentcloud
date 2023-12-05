Use this data source to query the detail information of redis instance.

Example Usage

```hcl
data "tencentcloud_redis_instances" "redislab" {
  zone               = "ap-hongkong-1"
  search_key         = "myredis"
  project_id         = 0
  limit              = 20
  result_output_file = "/tmp/redis_instances"
}
```