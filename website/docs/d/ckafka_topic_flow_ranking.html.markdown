---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_topic_flow_ranking"
sidebar_current: "docs-tencentcloud-datasource-ckafka_topic_flow_ranking"
description: |-
  Use this data source to query detailed information of ckafka topic_flow_ranking
---

# tencentcloud_ckafka_topic_flow_ranking

Use this data source to query detailed information of ckafka topic_flow_ranking

## Example Usage

```hcl
data "tencentcloud_ckafka_topic_flow_ranking" "topic_flow_ranking" {
  instance_id  = "ckafka-xxxxxx"
  ranking_type = "PRO"
  begin_date   = "2023-05-29T00:00:00+08:00"
  end_date     = "2021-05-29T23:59:59+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) InstanceId.
* `ranking_type` - (Required, String) Ranking type. `PRO`: topic production flow, `CON`: topic consumption traffic.
* `begin_date` - (Optional, String) BeginDate.
* `end_date` - (Optional, String) EndDate.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result.
  * `consume_speed` - ConsumeSpeed.
    * `consumer_group_name` - ConsumerGroupName.
    * `speed` - Speed.
  * `topic_flow` - TopicFlow.
    * `message_heap` - Topic MessageHeap.
    * `partition_num` - partitionNum.
    * `replica_num` - ReplicaNum.
    * `topic_id` - topicId.
    * `topic_name` - topicName.
    * `topic_traffic` - TopicTraffic.
  * `topic_message_heap` - TopicMessageHeapRanking.
    * `message_heap` - Topic MessageHeap.
    * `partition_num` - PartitionNum.
    * `replica_num` - ReplicaNum.
    * `topic_id` - topicId.
    * `topic_name` - topicName.
    * `topic_traffic` - TopicTraffic.


