---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_topic"
sidebar_current: "docs-tencentcloud-resource-ckafka_topic"
description: |-
  Use this resource to create ckafka topic.
---

# tencentcloud_ckafka_topic

Use this resource to create ckafka topic.

## Example Usage

```hcl
resource "tencentcloud_ckafka_topic" "example" {
  instance_id                    = "ckafka-bzmjpavn"
  topic_name                     = "tf-example"
  note                           = "topic note"
  replica_num                    = 4
  partition_num                  = 2
  enable_white_list              = true
  ip_white_list                  = ["1.1.1.1", "2.2.2.2"]
  clean_up_policy                = "delete"
  sync_replica_min_num           = 2
  unclean_leader_election_enable = false
  segment                        = 86400000
  retention                      = 60000
  max_message_bytes              = 4096
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Ckafka instance ID.
* `partition_num` - (Required, Int) The number of partition.
* `replica_num` - (Required, Int) The number of replica.
* `topic_name` - (Required, String, ForceNew) Name of the CKafka topic. It must start with a letter, the rest can contain letters, numbers and dashes(-).
* `clean_up_policy` - (Optional, String) Clear log policy, log clear mode, default is `delete`. `delete`: logs are deleted according to the storage time. `compact`: logs are compressed according to the key. `compact, delete`: logs are compressed according to the key and will be deleted according to the storage time.
* `enable_white_list` - (Optional, Bool) Whether to open the ip whitelist, `true`: open, `false`: close.
* `ip_white_list` - (Optional, List: [`String`]) Ip whitelist, quota limit, required when enableWhileList=true.
* `max_message_bytes` - (Optional, Int) Max message bytes. min: 1024 Byte(1KB), max: 8388608 Byte(8MB).
* `note` - (Optional, String) The subject note. It must start with a letter, and the remaining part can contain letters, numbers and dashes (-).
* `retention` - (Optional, Int) Message can be selected. Retention time, unit is ms, the current minimum value is 60000ms.
* `segment` - (Optional, Int) Segment scrolling time, in ms, the current minimum is 3600000ms.
* `sync_replica_min_num` - (Optional, Int) Min number of sync replicas, Default is `1`.
* `unclean_leader_election_enable` - (Optional, Bool) Whether to allow unsynchronized replicas to be selected as leader, default is `false`, `true: `allowed, `false`: not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the CKafka topic.
* `forward_cos_bucket` - Data backup cos bucket: the bucket address that is dumped to cos.
* `forward_interval` - Periodic frequency of data backup to cos.
* `forward_status` - Data backup cos status. Valid values: `0`, `1`. `1`: do not open data backup, `0`: open data backup.
* `message_storage_location` - Message storage location.
* `segment_bytes` - Number of bytes rolled by shard.


## Import

ckafka topic can be imported using the instance_id#topic_name, e.g.

```
$ terraform import tencentcloud_ckafka_topic.example ckafka-f9ife4zz#tf-example
```

