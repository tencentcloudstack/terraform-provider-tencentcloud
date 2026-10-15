---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_consumer_group"
sidebar_current: "docs-tencentcloud-resource-ckafka_consumer_group"
description: |-
  Provides a resource to create a ckafka consumer_group
---

# tencentcloud_ckafka_consumer_group

Provides a resource to create a ckafka consumer_group

## Example Usage

```hcl
resource "tencentcloud_ckafka_consumer_group" "consumer_group" {
  instance_id     = "InstanceId"
  group_name      = "GroupName"
  topic_name_list = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String, ForceNew) GroupName.
* `instance_id` - (Required, String, ForceNew) InstanceId.
* `topic_name_list` - (Optional, Set: [`String`], ForceNew) array of topic names.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateConsumer](https://cloud.tencent.com/document/api/597/70173) - 创建消费者组
* Read: [DescribeConsumerGroup](https://cloud.tencent.com/document/api/597/40841) - 查询消费分组信息
* Delete: [DeleteGroup](https://cloud.tencent.com/document/api/597/69059) - 删除消费组

## Import

ckafka consumer_group can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_consumer_group.consumer_group consumer_group_id
```

