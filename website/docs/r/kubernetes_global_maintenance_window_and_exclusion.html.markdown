---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_global_maintenance_window_and_exclusion"
sidebar_current: "docs-tencentcloud-resource-kubernetes_global_maintenance_window_and_exclusion"
description: |-
  Provides a resource to create a TKE kubernetes global maintenance window and exclusion
---

# tencentcloud_kubernetes_global_maintenance_window_and_exclusion

Provides a resource to create a TKE kubernetes global maintenance window and exclusion

~> **NOTE:** If the `target_regions` is `*`, global maintenance window with all region cannot be deleted

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `day_of_week` - (Required, Set: [`String`]) Maintenance cycle (which days of the week). supported parameter values are as follows:

- MO: Monday
- TU: Tuesday
- WE: Wednesday
- TH: Thursday
- FR: Friday
- SA: Saturday
- SU: Sunday.
* `duration` - (Required, Int) Maintenance duration (hours).
* `maintenance_time` - (Required, String) Maintenance start time.
* `target_regions` - (Required, Set: [`String`]) Regions.
* `exclusions` - (Optional, List) Maintenance exclusions.

The `exclusions` object supports the following:

* `end_at` - (Required, String) Maintenance exclusion end time.
* `name` - (Required, String) Maintenance exclusion name.
* `start_at` - (Required, String) Maintenance exclusion start time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TKE kubernetes global maintenance window and exclusion can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_global_maintenance_window_and_exclusion.example 55
```

