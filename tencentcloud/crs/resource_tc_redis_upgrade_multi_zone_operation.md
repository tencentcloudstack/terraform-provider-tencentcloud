Provides a resource to create a redis upgrade_multi_zone_operation

Example Usage

```hcl
resource "tencentcloud_redis_upgrade_multi_zone_operation" "upgrade_multi_zone_operation" {
  instance_id = "crs-c1nl9rpv"
  upgrade_proxy_and_redis_server = true
}
```