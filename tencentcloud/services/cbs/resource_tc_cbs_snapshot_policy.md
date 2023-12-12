Provides a snapshot policy resource.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "mysnapshotpolicyname"
  repeat_weekdays      = [1, 4]
  repeat_hours         = [1]
  retention_days       = 7
}
```

Import

CBS snapshot policy can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot_policy.snapshot_policy asp-jliex1tn
```