---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_group_offsets"
sidebar_current: "docs-tencentcloud-datasource-ckafka_group_offsets"
description: |-
  Use this data source to query detailed information of ckafka group_offsets
---

# tencentcloud_ckafka_group_offsets

Use this data source to query detailed information of ckafka group_offsets

## Example Usage

```hcl
data "tencentcloud_ckafka_group_offsets" "group_offsets" {
  instance_id = "ckafka-xxxxxx"
  group       = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String) Kafka consumer group name.
* `instance_id` - (Required, String) InstanceId.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) fuzzy match topicName.
* `topics` - (Optional, Set: [`String`]) An array of topic names subscribed by the group, if there is no such array, it means all topic information under the specified group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `topic_list` - The topic array, where each element is a json object.
  * `partitions` - he topic partition array, where each element is a json object.
    * `error_code` - ErrorCode.
    * `lag` - The number of unconsumed messages.
    * `log_end_offset` - The latest offset of the current partition.
    * `metadata` - When consumers submit messages, they can pass in metadata for other purposes. Currently, it is usually an empty string.
    * `offset` - The offset of the position.
    * `partition` - topic partitionId.
  * `topic` - topicName.


