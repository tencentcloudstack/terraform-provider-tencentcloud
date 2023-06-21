---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_messages"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rocketmq_messages"
description: |-
  Use this data source to query detailed information of tdmq message
---

# tencentcloud_tdmq_rocketmq_messages

Use this data source to query detailed information of tdmq message

## Example Usage

```hcl
data "tencentcloud_tdmq_rocketmq_messages" "message" {
  cluster_id     = "rocketmq-rkrbm52djmro"
  environment_id = "keep_ns"
  topic_name     = "keep-topic"
  msg_id         = "A9FE8D0567FE15DB97425FC08EEF0000"
  query_dlq_msg  = false
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster id.
* `environment_id` - (Required, String) Environment.
* `msg_id` - (Required, String) Message ID.
* `topic_name` - (Required, String) Topic, groupId is passed when querying dead letters.
* `query_dlq_msg` - (Optional, Bool) The value is true when querying dead letters, only valid for Rocketmq.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `body` - Message body.
* `message_tracks` - Consumer Group ConsumptionNote: This field may return null, indicating that no valid value can be obtained.
  * `consume_status` - consumption status.
  * `exception_desc` - Exception informationNote: This field may return null, indicating that no valid value can be obtained.
  * `group` - consumer group.
  * `track_type` - message track type.
* `produce_time` - Production time.
* `producer_addr` - Producer address.
* `properties` - Detailed parameters.
* `show_topic_name` - The topic name displayed on the details pageNote: This field may return null, indicating that no valid value can be obtained.


