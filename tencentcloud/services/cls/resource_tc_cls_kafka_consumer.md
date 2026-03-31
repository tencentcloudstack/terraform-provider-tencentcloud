Provides a resource to create a CLS kafka consumer

Example Usage

```hcl
resource "tencentcloud_cls_kafka_consumer" "example" {
  from_topic_id = "c9b68233-948a-4eaf-a363-d0c2ced393ae"
  compression   = 0
  consumer_content {
    enable_tag      = false
    format          = 1
    json_type       = 1
    meta_fields     = [
      "__SOURCE__",
      "__FILENAME__"
    ]
    tag_transaction = 2
  }
}
```

Import

CLS kafka consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_kafka_consumer.example c9b68233-948a-4eaf-a363-d0c2ced393ae
```
