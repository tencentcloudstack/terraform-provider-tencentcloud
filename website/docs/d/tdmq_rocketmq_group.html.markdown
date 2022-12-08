---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_group"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rocketmq_group"
description: |-
  Use this data source to query detailed information of tdmqRocketmq group
---

# tencentcloud_tdmq_rocketmq_group

Use this data source to query detailed information of tdmqRocketmq group

## Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  cluster_name = "test_rocketmq_datasource_group"
  remark       = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace_datasource"
  ttl            = 65000
  retention_time = 65000
  remark         = "test namespace"
}

resource "tencentcloud_tdmq_rocketmq_group" "group" {
  group_name       = "test_rocketmq_group"
  namespace        = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  read_enable      = true
  broadcast_enable = true
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  remark           = "test rocketmq group"
}

data "tencentcloud_tdmq_rocketmq_group" "group" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_id = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  filter_group = tencentcloud_tdmq_rocketmq_group.group.group_name
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace_id` - (Required, String) Namespace.
* `filter_group` - (Optional, String) Consumer group query by consumer group name. Fuzzy query is supported.
* `filter_one_group` - (Optional, String) Subscription group name. After it is specified, the information of only this subscription group will be returned.
* `filter_topic` - (Optional, String) Topic name, which can be used to query all subscription groups under the topic.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - List of subscription groups.
  * `broadcast_enable` - Whether to enable broadcast consumption.
  * `client_protocol` - Client protocol.
  * `consumer_num` - The number of online consumers.
  * `consumer_type` - Consumer type. Enumerated values: ACTIVELY or PASSIVELY.
  * `consumption_mode` - `0`: Cluster consumption mode; `1`: Broadcast consumption mode; `-1`: Unknown.
  * `create_time` - Creation time in milliseconds.
  * `name` - Consumer group name.
  * `read_enable` - Whether to enable consumption.
  * `remark` - Remarks (up to 128 characters).
  * `retry_partition_num` - The number of partitions in a retry topic.
  * `total_accumulative` - The total number of heaped messages.
  * `tps` - Consumption TPS.
  * `update_time` - Modification time in milliseconds.


