---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_topic"
sidebar_current: "docs-tencentcloud-resource-tdmq_rocketmq_topic"
description: |-
  Provides a resource to create a tdmqRocketmq topic
---

# tencentcloud_tdmq_rocketmq_topic

Provides a resource to create a tdmqRocketmq topic

## Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_topic" "example" {
  topic_name     = "tf_example"
  namespace_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  type           = "Normal"
  remark         = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace_name` - (Required, String) Topic namespace. Currently, you can create topics only in one single namespace.
* `topic_name` - (Required, String) Topic name, which can contain 3-64 letters, digits, hyphens, and underscores.
* `type` - (Required, String) Topic type. Valid values: Normal, GlobalOrder, PartitionedOrder.
* `partition_num` - (Optional, Int) Number of partitions.
* `remark` - (Optional, String) Topic remarks (up to 128 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in milliseconds.
* `update_time` - Update time in milliseconds.


## Import

tdmqRocketmq topic can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_topic.topic topic_id
```

