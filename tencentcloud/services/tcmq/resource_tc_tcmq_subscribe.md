Provides a resource to create a tcmq subscribe

Example Usage

```hcl
resource "tencentcloud_tcmq_subscribe" "subscribe" {
  topic_name = "topic_name"
  subscription_name = "subscription_name"
  protocol = "http"
  endpoint = "http://xxxxxx";
}
```

Import

tcmq subscribe can be imported using the id, e.g.

```
terraform import tencentcloud_tcmq_subscribe.subscribe subscribe_id
```