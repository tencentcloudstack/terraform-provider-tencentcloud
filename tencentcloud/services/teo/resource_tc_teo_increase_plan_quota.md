Provides a resource to increase TEO plan quota. Use this resource to purchase additional quota for a TEO enterprise plan when the number of bound sites, Web Protection custom precise match policy rules, or rate limiting precise rate limiting module rules reaches the plan's quota limit.

Example Usage

```hcl
resource "tencentcloud_teo_increase_plan_quota" "example" {
  plan_id      = "edgeone-2unuvzjmmn2q"
  quota_type   = "site"
  quota_number = 1
}
```