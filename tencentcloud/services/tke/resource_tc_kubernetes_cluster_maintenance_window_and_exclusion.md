Provides a resource to create a TKE kubernetes cluster maintenance window and exclusion

Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion" "example" {
  cluster_id       = "cls-d2cit6no"
  maintenance_time = "01:00:00"
  duration         = 4
  day_of_week      = ["MO", "TU", "WE", "TH", "FR"]
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

tke kubernetes_cluster_maintenance_window_and_exclusion can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.example cls-d2cit6no
```
