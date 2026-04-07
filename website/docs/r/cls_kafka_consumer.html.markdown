---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_kafka_consumer"
sidebar_current: "docs-tencentcloud-resource-cls_kafka_consumer"
description: |-
  Provides a resource to create a CLS kafka consumer
---

# tencentcloud_cls_kafka_consumer

Provides a resource to create a CLS kafka consumer

## Example Usage

```hcl
resource "tencentcloud_cls_kafka_consumer" "example" {
  from_topic_id = "c9b68233-948a-4eaf-a363-d0c2ced393ae"
  compression   = 0
  consumer_content {
    enable_tag = false
    format     = 1
    json_type  = 1
    meta_fields = [
      "__SOURCE__",
      "__FILENAME__"
    ]
    tag_transaction = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `from_topic_id` - (Required, String, ForceNew) Log topic ID.
* `compression` - (Optional, Int) Compression method: 0-NONE, 2-SNAPPY, 3-LZ4.
* `consumer_content` - (Optional, List) Kafka protocol consumption data format.

The `consumer_content` object supports the following:

* `enable_tag` - (Optional, Bool) Whether to deliver TAG information.
* `format` - (Optional, Int) Content format: 0-original content, 1-JSON.
* `json_type` - (Optional, Int) Consumption data JSON format: 1-not escaped, 2-escaped.
* `meta_fields` - (Optional, List) Metadata field list.
* `tag_transaction` - (Optional, Int) Tag data processing method: 1-not flattened, 2-flattened.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `topic_id` - Topic parameter used when KafkaConsumer consumes.


## Import

CLS kafka consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_kafka_consumer.example c9b68233-948a-4eaf-a363-d0c2ced393ae
```

