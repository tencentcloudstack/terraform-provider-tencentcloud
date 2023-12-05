Provide a resource to create a KubernetesClusterEndpoint. This resource allows you to create an empty cluster first without any workers. Only all attached node depends create complete, cluster endpoint will finally be enabled.

~> **NOTE:** Recommend using `depends_on` to make sure endpoint create after node pools or workers does.

Example Usage

```hcl
resource "tencentcloud_kubernetes_node_pool" "pool1" {}

resource "tencentcloud_kubernetes_cluster_endpoint" "foo" {
  cluster_id = "cls-xxxxxxxx"
  cluster_internet = true
  cluster_intranet = true
  # managed_cluster_internet_security_policies = [
    "192.168.0.0/24"
  ]
  cluster_intranet_subnet_id = "subnet-xxxxxxxx"
  depends_on = [
	tencentcloud_kubernetes_node_pool.pool1
  ]
}
```

Import

KubernetesClusterEndpoint instance can be imported by passing cluster id, e.g.
```
$ terraform import tencentcloud_kubernetes_cluster_endpoint.test cluster-id
```