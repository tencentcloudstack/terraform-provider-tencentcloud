Provides a resource to create a cynosdb cluster_resource_packages_attachment

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_resource_packages_attachment" "cluster_resource_packages_attachment" {
  cluster_id  = "cynosdbmysql-q1d8151n"
  package_ids = ["package-hy4d2ppl"]
}
```

Import

cynosdb cluster_resource_packages_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_resource_packages_attachment.cluster_resource_packages_attachment cluster_resource_packages_attachment_id
```