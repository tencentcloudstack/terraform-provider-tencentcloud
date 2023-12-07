Use this data source to query detailed information of redis instance_shards

Example Usage

```hcl
data "tencentcloud_redis_instance_shards" "instance_shards" {
  instance_id = "crs-c1nl9rpv"
  filter_slave = false
}
```