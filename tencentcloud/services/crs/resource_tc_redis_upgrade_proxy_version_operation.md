Provides a resource to create a redis upgrade proxy version operation

Example Usage

```hcl
resource "tencentcloud_redis_upgrade_proxy_version_operation" "example" {
  instance_id               = "crs-c1nl9rpv"
  current_proxy_version     = "5.0.0"
  upgrade_proxy_version     = "5.8.12"
  instance_type_upgrade_now = 1
}
```