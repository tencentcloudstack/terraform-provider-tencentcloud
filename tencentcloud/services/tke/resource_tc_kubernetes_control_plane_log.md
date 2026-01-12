Provides a resource to create a TKE kubernetes control plane log

Example Usage

Use automatic creation of log_set_id and topic_id

```hcl
resource "tencentcloud_kubernetes_control_plane_log" "example" {
  cluster_id   = "cls-rng1h5ei"
  cluster_type = "tke"
  components {
    name         = "karpenter"
    log_level    = "2"
    topic_region = "ap-guangzhou"
  }

  delete_log_set_and_topic = true
}
```

Use custom log_set_id and topic_id

```hcl
resource "tencentcloud_kubernetes_control_plane_log" "example" {
  cluster_id   = "cls-rng1h5ei"
  cluster_type = "tke"
  components {
    name         = "cluster-autoscaler"
    log_level    = "2"
    log_set_id   = "40eed846-0f43-44b1-b216-c786a8970b1f"
    topic_id     = "21918a54-9ab4-40bc-90cd-c600cff00695"
    topic_region = "ap-guangzhou"
  }
}
```

Import

TKE kubernetes control plane log can be imported using the clusterId#clusterType#componentName, e.g.

```
terraform import tencentcloud_kubernetes_control_plane_log.example cls-rng1h5ei#tke#cluster-autoscaler
```
