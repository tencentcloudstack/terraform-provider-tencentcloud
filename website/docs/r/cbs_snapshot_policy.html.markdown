---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot_policy"
sidebar_current: "docs-tencentcloud-resource-cbs_snapshot_policy"
description: |-
  Provides a CBS snapshot policy resource.
---

# tencentcloud_cbs_snapshot_policy

Provides a CBS snapshot policy resource.

~> **NOTE:** `repeat_weekdays`/`repeat_months`/`interval_days` are mutually exclusive rules, required and only one policy rule can be set. If mutually exclusive rules are passed in at the same time, only one will take effect, with priority in the order listed above: for example, if all three rules are set, only `repeat_weekdays` will take effect

## Example Usage

### Create policy by repeat_weekdays

```hcl
resource "tencentcloud_cbs_snapshot_policy" "example" {
  snapshot_policy_name = "tf-example"
  repeat_hours         = [1]
  repeat_weekdays      = [1, 4]
  is_activated         = true
  is_permanent         = false
  retention_days       = 7
}
```

### Create policy by interval_days

```hcl
resource "tencentcloud_cbs_snapshot_policy" "example" {
  snapshot_policy_name = "tf-example"
  repeat_hours         = [1]
  interval_days        = 30
  is_activated         = true
  is_permanent         = false
  retention_days       = 7
}
```

## Argument Reference

The following arguments are supported:

* `repeat_hours` - (Required, List: [`Int`]) Trigger times of periodic snapshot. Valid value ranges: (0~23). The 0 means 00:00, and so on.
* `snapshot_policy_name` - (Required, String) Name of snapshot policy. The maximum length can not exceed 60 bytes.
* `interval_days` - (Optional, Int) Specifies the interval in days for creating scheduled snapshots. Valid value ranges: [1, 365]. For example, setting it to 5 means a scheduled snapshot is created every 5 days. Note: when backing up by days, theoretically the first backup is on the day the policy is created; if the policy creation time is later than the configured backup time on that day, the first backup will be in the next backup cycle.
* `is_activated` - (Optional, Bool) Whether the snapshot policy is activated. Default is true.
* `is_permanent` - (Optional, Bool) Whether snapshots created by this policy are retained permanently. Default is false. Note: if this parameter is set to true, retention_days must not be specified, otherwise a conflict error will be returned by the API.
* `repeat_months` - (Optional, List: [`Int`]) Specifies the days of the month on which scheduled snapshots are triggered. Valid value ranges: [1, 31]. For example, 5 means triggering snapshots on the 5th day of each month. Note: if a date such as 29, 30 or 31 does not exist in some months, those months will be skipped.
* `repeat_weekdays` - (Optional, List: [`Int`]) Periodic snapshot is enabled. Valid values: [0, 1, 2, 3, 4, 5, 6]. 0 means Sunday, 1-6 means Monday to Saturday.
* `retention_days` - (Optional, Int) Retention days of the snapshot, and the default value is 7.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAutoSnapshotPolicy](https://www.tencentcloud.com/document/api/362/32167)
* Read: [DescribeAutoSnapshotPolicies](https://www.tencentcloud.com/document/api/362/32165)
* Update: [ModifyAutoSnapshotPolicyAttribute](https://www.tencentcloud.com/document/api/362/32162)
* Delete: [DeleteAutoSnapshotPolicies](https://www.tencentcloud.com/document/api/362/32166)

## Import

CBS snapshot policy can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_snapshot_policy.example asp-jliex1tn
```

