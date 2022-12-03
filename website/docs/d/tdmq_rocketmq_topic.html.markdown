---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_topic"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rocketmq_topic"
description: |-
  Use this data source to query detailed information of tdmqRocketmq topic
---

# tencentcloud_tdmq_rocketmq_topic

Use this data source to query detailed information of tdmqRocketmq topic

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

data "tencentcloud_tdmq_rocketmq_topic" "topic" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_id = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  filter_name  = tencentcloud_tdmq_rocketmq_topic.topic.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace_id` - (Required, String) Namespace.
* `filter_name` - (Optional, String) Search by topic name. Fuzzy query is supported.
* `filter_type` - (Optional, Set: [`String`]) Filter by topic type. Valid values: `Normal`, `GlobalOrder`, `PartitionedOrder`, `Transaction`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `topics` - List of topic information.
  * `create_time` - Creation time in milliseconds.
  * `name` - Topic name.
  * `partition_num` - The number of read/write partitions.
  * `remark` - Topic name.
  * `update_time` - Update time in milliseconds.


