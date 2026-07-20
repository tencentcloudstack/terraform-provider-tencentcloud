Provides a resource to purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. This is a one-time operation resource; the plan purchase is irreversible and cannot be cancelled on resource destroy.

Example Usage

```hcl
resource "tencentcloud_teo_plan_for_zone" "example" {
  zone_id   = "zone-27h0vbm5w1e"
  plan_type = "sta_global"
}
```

Purchase a standard plan with bot management

```hcl
resource "tencentcloud_teo_plan_for_zone" "example_bot" {
  zone_id   = "zone-27h0vbm5w1e"
  plan_type = "sta_global_with_bot"
}
```
