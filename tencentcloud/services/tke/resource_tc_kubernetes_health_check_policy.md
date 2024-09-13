Provides a resource to create a tke kubernetes_health_check_policy

Example Usage

```hcl
resource "tencentcloud_kubernetes_health_check_policy" "kubernetes_health_check_policy" {
    cluster_id = "cls-xxxxx"
    name = "example"
    rules {
        name = "OOMKilling"
        auto_repair_enabled = true
        enabled = true
    }
    rules {
        name = "KubeletUnhealthy"
        auto_repair_enabled = true
        enabled = true
    }
}
```

Import

tke kubernetes_health_check_policy can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy cls-xxxxx#healthcheckpolicyname
```
