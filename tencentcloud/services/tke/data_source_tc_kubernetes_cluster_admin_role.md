Provide a datasource to acquire TKE cluster admin role.

Use this data source to grant the current user (or sub-account) the `tke:admin` ClusterRole in the specified Kubernetes cluster. This is typically used when a CAM sub-account needs to be granted cluster administrator permissions through a CAM policy.

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_admin_role" "foo" {
  cluster_id = "cls-xxxxxxxx"
}

output "request_id" {
  value = data.tencentcloud_kubernetes_cluster_admin_role.foo.request_id
}
```
