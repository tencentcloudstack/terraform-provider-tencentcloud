---
subcategory: "Ckafka"
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
resource "tencentcloud_ckafka_topic" "foo" {
  instance_id                    = "ckafka-f9ife4zz"
  topic_name                     = "example"
  note                           = "topic note"
  replica_num                    = 2
  partition_num                  = 1
  enable_white_list              = true
  ip_white_list                  = ["ip1", "ip2"]
  clean_up_policy                = "delete"
  sync_replica_min_num           = 1
  unclean_leader_election_enable = false
  segment                        = 3600000
  retention                      = 60000
  max_message_bytes              = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Ckafka instance ID.
* `partition_num` - (Required) The number of partition.
* `replica_num` - (Required) The number of replica.
* `topic_name` - (Required, ForceNew) Name of the CKafka topic. It must start with a letter, the rest can contain letters, numbers and dashes(-).
* `clean_up_policy` - (Optional) Clear log policy, log clear mode, default is `delete`. `delete`: logs are deleted according to the storage time. `compact`: logs are compressed according to the key. `compact, delete`: logs are compressed according to the key and will be deleted according to the storage time.
* `enable_white_list` - (Optional) Whether to open the ip whitelist, `true`: open, `false`: close.
* `ip_white_list` - (Optional) Ip whitelist, quota limit, required when enableWhileList=true.
* `max_message_bytes` - (Optional) Max message bytes.
* `note` - (Optional) The subject note. It must start with a letter, and the remaining part can contain letters, numbers and dashes (-).
* `retention` - (Optional) Message can be selected. Retention time, unit is ms, the current minimum value is 60000ms.
* `segment` - (Optional) Segment scrolling time, in ms, the current minimum is 3600000ms.
* `sync_replica_min_num` - (Optional) Min number of sync replicas, Default is `1`.
* `unclean_leader_election_enable` - (Optional) Whether to allow unsynchronized replicas to be selected as leader, default is `false`, `true: `allowed, `false`: not allowed.

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
$ terraform import tencentcloud_ckafka_topic.foo ckafka-f9ife4zz#example
```

