Provides a resource to create a redis wan_address

Example Usage

```hcl
resource "tencentcloud_redis_wan_address" "wan_address" {
  instance_id = "crs-dekqpd8v"
}
```

Import

redis wan_address can be imported using the instanceId, e.g.

```
terraform import tencentcloud_redis_wan_address.wan_address crs-dekqpd8v
```