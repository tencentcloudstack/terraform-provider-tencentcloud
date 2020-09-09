---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot_policy"
sidebar_current: "docs-tencentcloud-resource-cbs_snapshot_policy"
description: |-
  Provides a snapshot policy resource.
---

# tencentcloud_cbs_snapshot_policy

Provides a snapshot policy resource.

## Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "mysnapshotpolicyname"
  repeat_weekdays      = [1, 4]
  repeat_hours         = [1]
  retention_days       = 7
}
```

## Argument Reference

The following arguments are supported:

* `repeat_hours` - (Required) Trigger times of periodic snapshot, the available values are 0 to 23. The 0 means 00:00, and so on.
* `repeat_weekdays` - (Required) Periodic snapshot is enabled, the available values are [0, 1, 2, 3, 4, 5, 6]. 0 means Sunday, 1-6 means Monday to Saturday.
* `snapshot_policy_name` - (Required) Name of snapshot policy. The maximum length can not exceed 60 bytes.
* `retention_days` - (Optional) Retention days of the snapshot, and the default value is 7.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CBS snapshot policy can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot_policy.snapshot_policy asp-jliex1tn
```

