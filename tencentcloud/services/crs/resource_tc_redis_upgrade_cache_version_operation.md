Provides a resource to create a redis upgrade cache version operation

Example Usage

```hcl
resource "tencentcloud_redis_upgrade_cache_version_operation" "example" {
  instance_id               = "crs-c1nl9rpv"
  current_redis_version     = "5.0.0"
  upgrade_redis_version     = "5.0.1"
  instance_type_upgrade_now = 1
}
```