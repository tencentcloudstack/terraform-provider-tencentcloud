Provides a CBS snapshot policy resource.

~> **NOTE:** `repeat_weekdays`/`repeat_months`/`interval_days` are mutually exclusive rules, required and only one policy rule can be set. If mutually exclusive rules are passed in at the same time, only one will take effect, with priority in the order listed above: for example, if all three rules are set, only `repeat_weekdays` will take effect

Example Usage

Create policy by repeat_weekdays

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

Create policy by interval_days

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

Import

CBS snapshot policy can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_snapshot_policy.example asp-jliex1tn
```