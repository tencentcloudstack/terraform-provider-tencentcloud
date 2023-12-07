Provides a resource to create a cls kafka_recharge

Example Usage

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
  tags                 = {
    "test" = "test",
  }
}

resource "tencentcloud_cls_kafka_recharge" "kafka_recharge" {
  topic_id = tencentcloud_cls_topic.topic.id
  name = "tf-example-recharge"
  kafka_type = 0
  offset = -2
  is_encryption_addr =true
  user_kafka_topics = "recharge"
  kafka_instance = "ckafka-qzoeaqx8"
  log_recharge_rule {
    recharge_type = "json_log"
    encoding_format = 0
    default_time_switch = true
  }
}

```

Import

cls kafka_recharge can be imported using the id, e.g.

```
terraform import tencentcloud_cls_kafka_recharge.kafka_recharge kafka_recharge_id
```