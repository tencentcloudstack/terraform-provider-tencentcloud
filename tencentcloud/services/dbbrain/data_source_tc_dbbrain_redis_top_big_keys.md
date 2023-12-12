Use this data source to query detailed information of dbbrain redis_top_big_keys

Example Usage

```hcl
data "tencentcloud_dbbrain_redis_top_big_keys" "redis_top_big_keys" {
	instance_id = local.redis_id
	date        = "%s"
	product     = "redis"
	sort_by     = "Capacity"
	key_type    = "string"
}
```