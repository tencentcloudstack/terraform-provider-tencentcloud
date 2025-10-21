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
data "tencentcloud_tdmq_rocketmq_topic" "example" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_id = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  filter_name  = tencentcloud_tdmq_rocketmq_topic.example.topic_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example"
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


