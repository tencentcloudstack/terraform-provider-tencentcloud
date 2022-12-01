---
subcategory: "RocketMQ"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_group"
sidebar_current: "docs-tencentcloud-resource-tdmq_rocketmq_group"
description: |-
  Provides a resource to create a tdmqRocketmq group
---

# tencentcloud_tdmq_rocketmq_group

Provides a resource to create a tdmqRocketmq group

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

resource "tencentcloud_tdmq_rocketmq_group" "group" {
  group_name       = "test_rocketmq_group"
  namespace        = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
  read_enable      = true
  broadcast_enable = true
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  remark           = "test rocketmq group"
}
```

## Argument Reference

The following arguments are supported:

* `broadcast_enable` - (Required, Bool) Whether to enable broadcast consumption.
* `cluster_id` - (Required, String) Cluster ID.
* `group_name` - (Required, String) Group name (8-64 characters).
* `namespace` - (Required, String) Namespace. Currently, only one namespace is supported.
* `read_enable` - (Required, Bool) Whether to enable consumption.
* `remark` - (Optional, String) Remarks (up to 128 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `client_protocol` - Client protocol.
* `consumer_num` - The number of online consumers.
* `consumer_type` - Consumer type. Enumerated values: ACTIVELY or PASSIVELY.
* `consumption_mode` - `0`: Cluster consumption mode; `1`: Broadcast consumption mode; `-1`: Unknown.
* `create_time` - Creation time in milliseconds.
* `retry_partition_num` - The number of partitions in a retry topic.
* `total_accumulative` - The total number of heaped messages.
* `tps` - Consumption TPS.
* `update_time` - Modification time in milliseconds.


## Import

tdmqRocketmq group can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_group.group group_id
```

