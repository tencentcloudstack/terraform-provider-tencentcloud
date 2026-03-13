Provides a resource to create a Monitor external cluster

Example Usage

```hcl
resource "tencentcloud_monitor_external_cluster" "example" {
  instance_id    = "prom-gzg3f1em"
  cluster_region = "ap-guangzhou"
  cluster_name   = "tf-external-cluster"

  external_labels {
    name  = "clusterName"
    value = "example"
  }

  external_labels {
    name  = "environment"
    value = "prod"
  }

  enable_external = false
}
```

Import

Monitor external cluster can be imported using the instanceId#clusterId, e.g.

```
terraform import tencentcloud_monitor_external_cluster.example prom-gzg3f1em#ecls-qi9v5opk
```