---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_topic_sync_replica"
sidebar_current: "docs-tencentcloud-datasource-ckafka_topic_sync_replica"
description: |-
  Use this data source to query detailed information of ckafka topic_sync_replica
---

# tencentcloud_ckafka_topic_sync_replica

Use this data source to query detailed information of ckafka topic_sync_replica

## Example Usage

```hcl
data "tencentcloud_ckafka_topic_sync_replica" "topic_sync_replica" {
  instance_id = "ckafka-xxxxxx"
  topic_name  = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) InstanceId.
* `topic_name` - (Required, String) TopicName.
* `out_of_sync_replica_only` - (Optional, Bool) Filter only unsynced replicas.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `topic_in_sync_replica_list` - Topic details and copy collection.
  * `begin_offset` - BeginOffset.
  * `end_offset` - EndOffset.
  * `in_sync_replica` - ISR.
  * `leader` - Leader Id.
  * `message_count` - Message Count.
  * `out_of_sync_replica` - Out Of Sync Replica.
  * `partition` - partition name.
  * `replica` - replica set.


