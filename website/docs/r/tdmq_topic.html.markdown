---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_topic"
sidebar_current: "docs-tencentcloud-resource-tdmq_topic"
description: |-
  Provides a resource to create a TDMQ topic.
---

# tencentcloud_tdmq_topic

Provides a resource to create a TDMQ topic.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_topic" "example" {
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 6
  pulsar_topic_type = 3
  remark            = "remark."
  tags = {
    "env"  = "test"
    "team" = "backend"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) The Dedicated Cluster Id.
* `environ_id` - (Required, String, ForceNew) The name of tdmq namespace.
* `partitions` - (Required, Int) The partitions of topic.
* `topic_name` - (Required, String, ForceNew) The name of topic to be created.
* `pulsar_topic_type` - (Optional, Int) Pulsar Topic Type 0: Non-persistent non-partitioned 1: Non-persistent partitioned 2: Persistent non-partitioned 3: Persistent partitioned.
* `remark` - (Optional, String) Description of the namespace.
* `tags` - (Optional, List) Tag description list.
* `topic_type` - (Optional, Int, **Deprecated**) This input will be gradually discarded and can be switched to PulsarTopicType parameter 0: Normal message; 1: Global sequential messages; 2: Local sequential messages; 3: Retrying queue; 4: Dead letter queue. The type of topic.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateTopic](https://cloud.tencent.com/document/api/1179/46088) - This API is used to add a message topic in the specified partition and type.
* Read: [DescribeTopics](https://cloud.tencent.com/document/api/1179/46086) - This API is used to get the list of topics under an environment.
* Update: [ModifyTopic](https://cloud.tencent.com/document/api/1179/46085) - This API is used to modify the topic remarks and number of partitions.
* Update: [TagResources](https://cloud.tencent.com/document/api/651/72280) - This API is used to create and bind a tag uniformly to multiple specified resources of multiple Tencent Cloud services.
* Update: [UnTagResources](https://cloud.tencent.com/document/api/651/72279) - This API is used to unbind a tag uniformly from multiple specified resources of multiple Tencent Cloud services.
* Delete: [DeleteTopics](https://cloud.tencent.com/document/api/1179/46087) - This API is used to batch delete topics.

