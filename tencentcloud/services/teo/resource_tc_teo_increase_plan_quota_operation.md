Provides a resource to increase plan quota for TEO (EdgeOne) plans.

Example Usage

```hcl
resource "tencentcloud_teo_increase_plan_quota_operation" "example" {
  plan_id      = "edgeone-2unuvzjmmn2q"
  quota_type   = "site"
  quota_number = 10
}
```