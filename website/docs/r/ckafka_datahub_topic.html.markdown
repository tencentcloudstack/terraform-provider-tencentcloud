---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_datahub_topic"
sidebar_current: "docs-tencentcloud-resource-ckafka_datahub_topic"
description: |-
  Provides a resource to create a ckafka datahub_topic
---

# tencentcloud_ckafka_datahub_topic

Provides a resource to create a ckafka datahub_topic

## Example Usage

```hcl
data "tencentcloud_user_info" "user" {}

resource "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
  name          = format("%s-tf", data.tencentcloud_user_info.user.app_id)
  partition_num = 20
  retention_ms  = 60000
  note          = "for test"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Name, start with appid, which is a string of no more than 128 characters, must start with a letter, and the rest can contain letters, numbers, and dashes (-).
* `partition_num` - (Required, Int) Number of Partitions, greater than 0.
* `retention_ms` - (Required, Int) Message retention time, in ms, the current minimum value is 60000 ms.
* `note` - (Optional, String) Subject note, which is a string of no more than 64 characters, must start with a letter, and the rest can contain letters, numbers and dashes (-).
* `tags` - (Optional, Map) Tags of dataHub topic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDatahubTopic](https://cloud.tencent.com/document/api/597/89017) - This API is used to create a DIP topic.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeDatahubTopic](https://cloud.tencent.com/document/api/597/86862) - This API is used to retrieve DIP topic attributes.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyDatahubTopic](https://cloud.tencent.com/document/api/597/89542) - This API is used to modify DIP topic attributes.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteDatahubTopic](https://cloud.tencent.com/document/api/597/77803)

## Import

ckafka datahub_topic can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_datahub_topic.datahub_topic datahub_topic_name
```

