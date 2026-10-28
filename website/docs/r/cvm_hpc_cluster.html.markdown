---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_hpc_cluster"
sidebar_current: "docs-tencentcloud-resource-cvm_hpc_cluster"
description: |-
  Provides a resource to create a cvm hpc_cluster
---

# tencentcloud_cvm_hpc_cluster

Provides a resource to create a cvm hpc_cluster

## Example Usage

```hcl
resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone   = "ap-beijing-6"
  name   = "terraform-test"
  remark = "create for test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of Hpc Cluster.
* `zone` - (Required, String) Available zone.
* `remark` - (Optional, String) Remark of Hpc Cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateHpcCluster](https://www.tencentcloud.com/search/?wd=CreateHpcCluster&type=document&lang=en)
* Read: [DescribeHpcClusters](https://www.tencentcloud.com/search/?wd=DescribeHpcClusters&type=document&lang=en)
* Update: [ModifyHpcClusterAttribute](https://www.tencentcloud.com/search/?wd=ModifyHpcClusterAttribute&type=document&lang=en)
* Delete: [DeleteHpcClusters](https://www.tencentcloud.com/search/?wd=DeleteHpcClusters&type=document&lang=en)

## Import

cvm hpc_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_hpc_cluster.hpc_cluster hpc_cluster_id
```

