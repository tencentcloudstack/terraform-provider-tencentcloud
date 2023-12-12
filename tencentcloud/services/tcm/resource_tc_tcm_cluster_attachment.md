Provides a resource to create a tcm cluster_attachment

Example Usage

```hcl
resource "tencentcloud_tcm_cluster_attachment" "cluster_attachment" {
  mesh_id = "mesh-b9q6vf9l"
  cluster_list {
    cluster_id = "cls-rc5uy6dy"
    region = "ap-guangzhou"
    role = "REMOTE"
    vpc_id = "vpc-a1jycmbx"
    subnet_id = "subnet-lkyb3ayc"
    type = "TKE"
  }
}

```
Import

tcm cluster_attachment can be imported using the mesh_id#cluster_id, e.g.
```
$ terraform import tencentcloud_tcm_cluster_attachment.cluster_attachment mesh-b9q6vf9l#cls-rc5uy6dy
```