Provides a resource to create a tcmq topic

Example Usage

```hcl
resource "tencentcloud_tcmq_topic" "topic" {
  topic_name = "topic_name"
}
```

Import

tcmq topic can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_topic.topic topic_id
```