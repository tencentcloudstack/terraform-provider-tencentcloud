Provides a resource to create a tke kubernetes_health_check_policy

Example Usage

```hcl
resource "tencentcloud_kubernetes_health_check_policy" "kubernetes_health_check_policy" {
  health_check_policy = {
    rules = {
    }
  }
}
```

Import

tke kubernetes_health_check_policy can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy kubernetes_health_check_policy_id
```
