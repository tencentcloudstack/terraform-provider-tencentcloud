Provides a resource to create a cfs auto_snapshot_policy

Example Usage

use day of week

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  day_of_week = "1,2"
  hour = "2,3"
  policy_name = "policy_name"
  alive_days = 7
}
```

use day of month

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  hour = "2,3"
  policy_name = "policy_name"
  alive_days = 7
  day_of_month = "2,3,4"
}
```

use interval days

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  hour = "2,3"
  policy_name = "policy_name"
  alive_days = 7
  interval_days = 1
}
```

Import

cfs auto_snapshot_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_auto_snapshot_policy.auto_snapshot_policy auto_snapshot_policy_id
```