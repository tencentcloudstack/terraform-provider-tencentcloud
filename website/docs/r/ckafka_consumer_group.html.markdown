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

* Create: [CreateConsumer](https://www.tencentcloud.com/document/product/597/45903) - This API is used to create a consumer group.
* Read: [DescribeConsumerGroup](https://www.tencentcloud.com/document/product/597/35351) - This API is used to query consumer group information.
* Delete: [DeleteGroup](https://www.tencentcloud.com/document/product/597/74427)

## Import

ckafka consumer_group can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_consumer_group.consumer_group consumer_group_id
```

