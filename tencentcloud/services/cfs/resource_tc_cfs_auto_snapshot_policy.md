Provides a resource to create a cfs auto snapshot policy

Example Usage

Use day of week

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "example" {
  policy_name = "tf-example"
  day_of_week = "1,2"
  hour        = "2,3"
  alive_days  = 7
}
```

Use day of month

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "example" {
  policy_name  = "tf-example"
  day_of_month = "2,3,4"
  hour         = "2,3"
  alive_days   = 7
}
```

Use interval days

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy" "example" {
  policy_name   = "policy_name"
  interval_days = 1
  hour          = "2,3"
  alive_days    = 7
}
```

Import

cfs auto snapshot policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_auto_snapshot_policy.example asp-f8q793kj
```