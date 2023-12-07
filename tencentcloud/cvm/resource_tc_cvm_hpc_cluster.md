Provides a resource to create a cvm hpc_cluster

Example Usage

```hcl
resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone = "ap-beijing-6"
  name = "terraform-test"
  remark = "create for test"
}
```

Import

cvm hpc_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_hpc_cluster.hpc_cluster hpc_cluster_id
```