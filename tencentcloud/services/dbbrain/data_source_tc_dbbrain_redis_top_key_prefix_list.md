Use this data source to query detailed information of dbbrain redis_top_key_prefix_list

Example Usage

```hcl
data "tencentcloud_dbbrain_redis_top_key_prefix_list" "redis_top_key_prefix_list" {
	instance_id = local.redis_id
	date        = "%s"
	product     = "redis"
}
```