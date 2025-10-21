---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_topics"
sidebar_current: "docs-tencentcloud-datasource-ckafka_topics"
description: |-
  Use this data source to query detailed information of ckafka topic.
---

# tencentcloud_ckafka_topics

Use this data source to query detailed information of ckafka topic.

## Example Usage

```hcl
data "tencentcloud_ckafka_topics" "example" {
  instance_id = "ckafka-vv7wp5nx"
  topic_name  = "tf_example"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Ckafka instance ID.
* `result_output_file` - (Optional, String) Used to store results.
* `topic_name` - (Optional, String) Name of the CKafka topic. It must start with a letter, the rest can contain letters, numbers and dashes(-). The length range is from 1 to 64.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of instances. Each element contains the following attributes.
  * `clean_up_policy` - Clear log policy, log clear mode. `delete`: logs are deleted according to the storage time, `compact`: logs are compressed according to the key, `compact, delete`: logs are compressed according to the key and will be deleted according to the storage time.
  * `create_time` - Create time of the CKafka topic.
  * `enable_white_list` - Whether to open the IP Whitelist. `true`: open, `false`: close.
  * `forward_cos_bucket` - Data backup cos bucket: the bucket address that is dumped to cos.
  * `forward_interval` - Periodic frequency of data backup to cos.
  * `forward_status` - Data backup cos status. `1`: do not open data backup, `0`: open data backup.
  * `ip_white_list_count` - IP Whitelist count.
  * `max_message_bytes` - Max message bytes.
  * `note` - CKafka topic note description.
  * `partition_num` - The number of partition.
  * `replica_num` - The number of replica.
  * `retention` - Message can be selected. Retention time(unit ms).
  * `segment_bytes` - Number of bytes rolled by shard.
  * `segment` - Segment scrolling time, in ms.
  * `sync_replica_min_num` - Min number of sync replicas.
  * `topic_id` - ID of the CKafka topic.
  * `topic_name` - Name of the CKafka topic.
  * `unclean_leader_election_enable` - Whether to allow unsynchronized replicas to be selected as leader, default is `false`, `true: `allowed, `false`: not allowed.


