Provides a resource to create a TEO plan

Example Usage

```hcl
resource "tencentcloud_teo_plan" "example" {
  plan_type = "standard"
  prepaid_plan_param {
    period     = 1
    renew_flag = "on"
  }
}
```

Import

TEO plan can be imported using the id, e.g.

```
terraform import tencentcloud_teo_plan.example edgeone-3dnkntfojm6o
```
