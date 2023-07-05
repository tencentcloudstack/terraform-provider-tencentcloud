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

### Add serverless node pool to a cluster

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  sg_id     = data.tencentcloud_security_groups.sg.security_groups.0.security_group_id
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = local.vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "tf example cluster"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"
  cluster_deploy_type     = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_serverless_node_pool" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  name       = "tf_example_serverless_node_pool"

  serverless_nodes {
    display_name = "tf_example_serverless_node1"
    subnet_id    = local.subnet_id
  }

  serverless_nodes {
    display_name = "tf_example_serverless_node2"
    subnet_id    = local.subnet_id
  }

  security_group_ids = [local.sg_id]
  labels = {
    "label1" : "value1",
    "label2" : "value2",
  }
}
```

### Adding taints to the virtual nodes under this node pool

The pods without appropriate tolerations will not be scheduled on this node. Refer [taint-and-toleration](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) for more details.

```hcl
resource "tencentcloud_kubernetes_serverless_node_pool" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  name       = "tf_example_serverless_node_pool"

  serverless_nodes {
    display_name = "tf_example_serverless_node1"
    subnet_id    = local.subnet_id
  }

  serverless_nodes {
    display_name = "tf_example_serverless_node2"
    subnet_id    = local.subnet_id
  }

  security_group_ids = [local.sg_id]
  labels = {
    "label1" : "value1",
    "label2" : "value2",
  }

  taints {
    key    = "key1"
    value  = "value1"
    effect = "NoSchedule"
  }

  taints {
    key    = "key1"
    value  = "value1"
    effect = "NoExecute"
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

