Provides a resource to purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. This is a one-time operation action; the plan purchase is irreversible and cannot be cancelled on resource destroy.

Example Usage

Option 1: Invoke on demand via CLI

```sh
terraform apply -invoke=action.tencentcloud_teo_plan_for_zone.example
```

```hcl
action "tencentcloud_teo_plan_for_zone" "example" {
  config {
    zone_id   = "zone-27h0vbm5w1e"
    plan_type = "sta_global"
  }
}
```

Option 2: Triggered by the zone resource lifecycle

```hcl
resource "tencentcloud_teo_zone" "example" {
  # ...
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.tencentcloud_teo_plan_for_zone.example]
    }
  }
}
```
