Provides a resource to create a tke kubernetes_cluster_master_attachment

Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_master_attachment" "kubernetes_cluster_master_attachment" {
  extra_args = {
  }
  master_config = {
    labels = {
    }
    data_disks = {
    }
    extra_args = {
    }
    gpu_args = {
    }
    taints = {
    }
  }
}
```

Import

tke kubernetes_cluster_master_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_cluster_master_attachment.kubernetes_cluster_master_attachment kubernetes_cluster_master_attachment_id
```
