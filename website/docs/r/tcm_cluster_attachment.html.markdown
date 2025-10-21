---
subcategory: "TencentCloud ServiceMesh(TCM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcm_cluster_attachment"
sidebar_current: "docs-tencentcloud-resource-tcm_cluster_attachment"
description: |-
  Provides a resource to create a tcm cluster_attachment
---

# tencentcloud_tcm_cluster_attachment

Provides a resource to create a tcm cluster_attachment

## Example Usage

```hcl
resource "tencentcloud_tcm_cluster_attachment" "cluster_attachment" {
  mesh_id = "mesh-b9q6vf9l"
  cluster_list {
    cluster_id = "cls-rc5uy6dy"
    region     = "ap-guangzhou"
    role       = "REMOTE"
    vpc_id     = "vpc-a1jycmbx"
    subnet_id  = "subnet-lkyb3ayc"
    type       = "TKE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `mesh_id` - (Required, String, ForceNew) Mesh ID.
* `cluster_list` - (Optional, List, ForceNew) Cluster list.

The `cluster_list` object supports the following:

* `cluster_id` - (Required, String) TKE Cluster id.
* `region` - (Required, String) TKE cluster region.
* `role` - (Required, String) Cluster role in mesh, REMOTE or MASTER.
* `type` - (Required, String) Cluster type.
* `vpc_id` - (Required, String) Cluster&#39;s VpcId.
* `subnet_id` - (Optional, String) Subnet id, only needed if it&#39;s standalone mesh.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tcm cluster_attachment can be imported using the mesh_id#cluster_id, e.g.
```
$ terraform import tencentcloud_tcm_cluster_attachment.cluster_attachment mesh-b9q6vf9l#cls-rc5uy6dy
```

