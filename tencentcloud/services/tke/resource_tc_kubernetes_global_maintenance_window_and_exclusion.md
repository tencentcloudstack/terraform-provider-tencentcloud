Provides a resource to create a TKE kubernetes global maintenance window and exclusion

~> **NOTE:** If the `target_regions` is `*`, global maintenance window with all region cannot be deleted

Example Usage

```hcl
resource "tencentcloud_kubernetes_global_maintenance_window_and_exclusion" "example" {
  maintenance_time = "02:00:00"
  duration         = 4
  day_of_week      = ["MO", "TU", "WE", "TH", "FR"]
  target_regions   = ["ap-guangzhou"]
  exclusions {
    name     = "name1"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }

  exclusions {
    name     = "name2"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }

  exclusions {
    name     = "name3"
    start_at = "2026-03-01 23:59:59"
    end_at   = "2026-03-07 23:59:59"
  }
}
```

Import

TKE kubernetes global maintenance window and exclusion can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_global_maintenance_window_and_exclusion.example 55
```
