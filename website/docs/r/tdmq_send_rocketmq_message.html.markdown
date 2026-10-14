---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_send_rocketmq_message"
sidebar_current: "docs-tencentcloud-resource-tdmq_send_rocketmq_message"
description: |-
  Provides a resource to create a tdmq send_rocketmq_message
---

# tencentcloud_tdmq_send_rocketmq_message

Provides a resource to create a tdmq send_rocketmq_message

## Example Usage

```hcl
resource "tencentcloud_tdmq_send_rocketmq_message" "send_rocketmq_message" {
  cluster_id   = "rocketmq-7k45z9dkpnne"
  namespace_id = "test_ns"
  topic_name   = "test_topic"
  msg_body     = "msg key"
  msg_key      = "msg tag"
  msg_tag      = "msg value"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster id.
* `msg_body` - (Required, String, ForceNew) Information.
* `namespace_id` - (Required, String, ForceNew) Namespaces.
* `topic_name` - (Required, String, ForceNew) topic name.
* `msg_key` - (Optional, String, ForceNew) Message key information.
* `msg_tag` - (Optional, String, ForceNew) Message tag information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [SendRocketMQMessage](https://cloud.tencent.com/document/api/1179/94179) - This API is used to send messages through RocketMQ. It is only used for sending a small number of test messages from the console and does not provide SLA. Cloud API is subject to traffic throttling. In real business scenarios, use the RocketMQ SDK to send messages.

