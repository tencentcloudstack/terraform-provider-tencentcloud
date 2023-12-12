Provides a resource to create a redis upgrade_proxy_version_operation

Example Usage

```hcl
resource "tencentcloud_redis_upgrade_proxy_version_operation" "upgrade_proxy_version_operation" {
  instance_id = "crs-c1nl9rpv"
  current_proxy_version = "5.0.0"
  upgrade_proxy_version = "5.0.0"
  instance_type_upgrade_now = 1
}
```