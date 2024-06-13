Provide a resource to configure kubernetes cluster app addons.

Example Usage

Install cos addon

```hcl

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = "vpc-xxxxxxxx"
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_internet        = false # (can be ignored) open it after the nodes added
  cluster_version         = "1.22.5"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  # without any worker config
}

resource "tencentcloud_kubernetes_addon" "kubernetes_addon" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  addon_name    = "cos"
  addon_version = "2018-05-25"
  raw_values    = "e30="
}

```

Import

Addon can be imported by using cluster_id#addon_name
```
$ terraform import tencentcloud_kubernetes_addon.addon_cos cls-xxx#addon_name
```