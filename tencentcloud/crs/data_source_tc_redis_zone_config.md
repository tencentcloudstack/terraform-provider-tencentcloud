Use this data source to query which instance types of Redis are available in a specific region.

Example Usage

```hcl
data "tencentcloud_redis_zone_config" "redislab" {
  region             = "ap-hongkong"
  result_output_file = "/temp/mytestpath"
}
```