---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_kafka_recharge"
sidebar_current: "docs-tencentcloud-resource-cls_kafka_recharge"
description: |-
  Provides a resource to create a cls kafka_recharge
---

# tencentcloud_cls_kafka_recharge

Provides a resource to create a cls kafka_recharge

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-example-logset"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "tf-example-topic"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    "test" = "test",
  }
}

resource "tencentcloud_cls_kafka_recharge" "kafka_recharge" {
  topic_id           = tencentcloud_cls_topic.topic.id
  name               = "tf-example-recharge"
  kafka_type         = 0
  offset             = -2
  is_encryption_addr = true
  user_kafka_topics  = "recharge"
  kafka_instance     = "ckafka-qzoeaqx8"
  log_recharge_rule {
    recharge_type       = "json_log"
    encoding_format     = 0
    default_time_switch = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `kafka_type` - (Required, Int) kafka recharge type, 0 for CKafka, 1 fro user define Kafka.
* `name` - (Required, String) kafka recharge name.
* `offset` - (Required, Int) The translation is: -2: Earliest (default) -1: Latest.
* `topic_id` - (Required, String) recharge for cls TopicId.
* `user_kafka_topics` - (Required, String) user need recharge kafka topic list.
* `consumer_group_name` - (Optional, String) user consumer group name.
* `is_encryption_addr` - (Optional, Bool) ServerAddr is encryption addr.
* `kafka_instance` - (Optional, String) CKafka Instance id.
* `log_recharge_rule` - (Optional, List) log recharge rule.
* `protocol` - (Optional, List) encryption protocol.
* `server_addr` - (Optional, String) Server addr.

The `log_recharge_rule` object supports the following:

* `default_time_switch` - (Required, Bool) user default time.
* `encoding_format` - (Required, Int) encoding format.
* `recharge_type` - (Required, String) recharge type.
* `default_time_src` - (Optional, Int) default time from.
* `keys` - (Optional, Set) log key list.
* `log_regex` - (Optional, String) log regex.
* `metadata` - (Optional, Set) metadata.
* `time_format` - (Optional, String) time format.
* `time_key` - (Optional, String) time key.
* `time_regex` - (Optional, String) time regex.
* `time_zone` - (Optional, String) time zone.
* `un_match_log_key` - (Optional, String) parse failed log key.
* `un_match_log_switch` - (Optional, Bool) is push parse failed log.
* `un_match_log_time_src` - (Optional, Int) parse failed log time from.

The `protocol` object supports the following:

* `mechanism` - (Optional, String) encryption type.
* `password` - (Optional, String) user password.
* `protocol` - (Optional, String) protocol type.
* `user_name` - (Optional, String) username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls kafka_recharge can be imported using the id, e.g.

```
terraform import tencentcloud_cls_kafka_recharge.kafka_recharge kafka_recharge_id
```

