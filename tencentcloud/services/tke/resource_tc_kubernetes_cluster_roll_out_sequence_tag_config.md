Provides a resource to create a TKE kubernetes cluster roll out sequence tag config.

Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config" "example" {
  cluster_id = "cls-6ml650mu"

  tags {
    key   = "Env"
    value = "Test"
  }

  tags {
    key   = "Protection-Level"
    value = "Medium"
  }
}
```

Import

TKE kubernetes cluster roll out sequence tag config can be imported using the cluster_id, e.g.

```
terraform import tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example cls-6ml650mu
```
