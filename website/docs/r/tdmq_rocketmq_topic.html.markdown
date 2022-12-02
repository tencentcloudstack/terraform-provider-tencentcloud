---
subcategory: "RocketMQ"
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
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  cluster_name = "test_rocketmq"
  remark       = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace"
  ttl            = 65000
  retention_time = 65000
  remark         = "test namespace"
}

resource "tencentcloud_tdmq_rocketmq_topic" "topic" {
  topic_name     = "test_rocketmq_topic"
  namespace_name = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  type           = "Normal"
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  remark         = "test rocketmq topic"
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

