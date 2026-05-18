Provides a resource to manage TKE cluster extra args configuration.

~> **NOTE:** This resource must exclusive in one cluster, do not declare additional args resources of this extra args elsewhere.

Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_extra_args_config" "example" {
  cluster_id = "cls-man1vvi2"
  kube_apiserver = [
    "goaway-chance=0",
    "kubelet-preferred-address-types=Hostname"
  ]

  kube_controller_manager = [
    "concurrent-serviceaccount-token-syncs=5"
  ]

  kube_scheduler = [
    "kube-api-qps=50"
  ]
}
```

Import

TKE cluster extra args config can be imported using the clusterId, e.g.

```
terraform import tencentcloud_kubernetes_cluster_extra_args_config.example cls-man1vvi2
```
