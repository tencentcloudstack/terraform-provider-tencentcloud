Use this action to increase plan quota for TEO (EdgeOne) plans.

This action calls the TEO `IncreasePlanQuota` API to purchase additional plan quotas (e.g., site count, precise access control rules, rate limiting rules). It is a one-time operation that does not persist any state.

Example Usage

```hcl
action "tencentcloud_teo_increase_plan_quota" "example" {
  plan_id      = "edgeone-2unuvzjmmn2q"
  quota_type   = "site"
  quota_number = 10
}
```