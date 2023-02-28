---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_serverless_node_pool"
sidebar_current: "docs-tencentcloud-resource-kubernetes_serverless_node_pool"
description: |-
  Provide a resource to create serverless node pool of cluster.
---

# tencentcloud_kubernetes_serverless_node_pool

Provide a resource to create serverless node pool of cluster.

## Example Usage

```hcl
resource "tencentcloud_kubernetes_serverless_node_pool" "example_serverless_node_pool" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  name       = "example_node_pool"
  serverless_nodes {
    display_name = "serverless_node1"
    subnet_id    = "subnet-xxx"
  }
  serverless_nodes {
    display_name = "serverless_node2"
    subnet_id    = "subnet-xxx"
  }
  security_group_ids = ["sg-xxx"]
  labels = {
    "example1" : "test1",
    "example2" : "test2",
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) cluster id of serverless node pool.
* `serverless_nodes` - (Required, List, ForceNew) node list of serverless node pool.
* `labels` - (Optional, Map) labels of serverless node.
* `name` - (Optional, String) serverless node pool name.
* `security_group_ids` - (Optional, List: [`String`], ForceNew) security groups of serverless node pool.
* `taints` - (Optional, List) taints of serverless node.

The `serverless_nodes` object supports the following:

* `subnet_id` - (Required, String) subnet id of serverless node.
* `display_name` - (Optional, String) display name of serverless node.

The `taints` object supports the following:

* `effect` - (Required, String) Effect of the taint. Valid values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.
* `key` - (Required, String) Key of the taint. The taint key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').
* `value` - (Required, String) Value of the taint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `life_state` - life state of serverless node pool.


## Import

serverless node pool can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_serverless_node_pool.test cls-xxx#np-xxx
```

