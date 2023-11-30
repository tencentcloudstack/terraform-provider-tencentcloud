Provides a resource to create a CLB instance topic.

Example Usage

```hcl
resource "tencentcloud_clb_log_topic" "topic" {
  log_set_id = "${tencentcloud_clb_log_set.set.id}"
  topic_name = "clb-topic"
}
```

Import

CLB log topic can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_log_topic.topic lb-7a0t6zqb