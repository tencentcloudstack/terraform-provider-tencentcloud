Provides a resource to create a tdmq send_rocketmq_message

Example Usage

```hcl
resource "tencentcloud_tdmq_send_rocketmq_message" "send_rocketmq_message" {
  cluster_id   = "rocketmq-7k45z9dkpnne"
  namespace_id = "test_ns"
  topic_name   = "test_topic"
  msg_body     = "msg key"
  msg_key      = "msg tag"
  msg_tag      = "msg value"
}
```