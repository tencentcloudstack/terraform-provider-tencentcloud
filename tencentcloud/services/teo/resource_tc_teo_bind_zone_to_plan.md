Provides a resource to bind a TEO zone to an existing plan. After creating an EdgeOne (TEO) zone, you can use this resource to bind the unbound zone to an existing plan so that the zone takes effect.

Example Usage

```hcl
resource "tencentcloud_teo_bind_zone_to_plan" "example" {
  zone_id = "zone-12345678"
  plan_id = "edgeone-12345678"
}
```
