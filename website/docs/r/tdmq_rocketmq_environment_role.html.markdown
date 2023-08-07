---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_environment_role"
sidebar_current: "docs-tencentcloud-resource-tdmq_rocketmq_environment_role"
description: |-
  Provides a resource to create a tdmqRocketmq environment_role
---

# tencentcloud_tdmq_rocketmq_environment_role

Provides a resource to create a tdmqRocketmq environment_role

## Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  role_name  = "tf_example_role"
  remark     = "remark."
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_environment_role" "example" {
  environment_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  role_name        = tencentcloud_tdmq_rocketmq_role.example.role_name
  permissions      = ["produce", "consume"]
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID (required).
* `environment_name` - (Required, String) Environment (namespace) name.
* `permissions` - (Required, Set: [`String`]) Permissions, which is a non-empty string array of `produce` and `consume` at the most.
* `role_name` - (Required, String) Role Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmqRocketmq environment_role can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_environment_role.environment_role environmentRole_id
```

