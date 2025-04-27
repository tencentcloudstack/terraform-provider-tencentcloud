Provides a resource to create a redis upgrade multi zone operation

Example Usage

```hcl
resource "tencentcloud_redis_upgrade_multi_zone_operation" "example" {
  instance_id                    = "crs-c1nl9rpv"
  upgrade_proxy_and_redis_server = true
}
```