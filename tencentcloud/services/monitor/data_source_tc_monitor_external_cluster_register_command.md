Use this data source to query Monitor external cluster register command

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

data "tencentcloud_monitor_external_cluster_register_command" "example" {
  instance_id = tencentcloud_monitor_external_cluster.example.instance_id
  cluster_id  = tencentcloud_monitor_external_cluster.example.cluster_id
}
```
