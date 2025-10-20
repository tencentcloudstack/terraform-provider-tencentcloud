Provides a resource to create a TKE kubernetes health check policy

Example Usage

```hcl
resource "tencentcloud_kubernetes_health_check_policy" "example" {
  cluster_id = "cls-fdy7hm1q"
  name       = "tf-example"
  rules {
    name                = "OOMKilling"
    auto_repair_enabled = true
    enabled             = true
  }

  rules {
    name                = "KubeletUnhealthy"
    auto_repair_enabled = true
    enabled             = true
  }
}
```

Import

TKE kubernetes health check policy can be imported using the clusterId#name, e.g.

```
terraform import tencentcloud_kubernetes_health_check_policy.example cls-fdy7hm1q#tf-example
```
