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

* Create: [CreateHpcCluster](https://cloud.tencent.com/document/api/213/83222) - 创建高性能计算集群
* Read: [DescribeHpcClusters](https://cloud.tencent.com/document/api/213/83220) - 查询高性能集群信息
* Update: [ModifyHpcClusterAttribute](https://cloud.tencent.com/document/api/213/83219) - 修改高性能计算集群属性
* Delete: [DeleteHpcClusters](https://cloud.tencent.com/document/api/213/83221) - 删除高性能计算集群

## Import

cvm hpc_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_hpc_cluster.hpc_cluster hpc_cluster_id
```

