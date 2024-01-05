Provides a resource to create a tdmq subscription

Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
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
  cluster_id        = tencentcloud_tdmq_instance.example.id
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  topic_name        = "tf-example-topic"
  partitions        = 1
  pulsar_topic_type = 3
  remark            = "remark."
}

resource "tencentcloud_tdmq_subscription" "example" {
  cluster_id               = tencentcloud_tdmq_instance.example.id
  environment_id           = tencentcloud_tdmq_namespace.example.environ_name
  topic_name               = tencentcloud_tdmq_topic.example.topic_name
  subscription_name        = "tf-example-subscription"
  remark                   = "remark."
  auto_create_policy_topic = true
  auto_delete_policy_topic = true
}
```

Import

tdmq subscription can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_subscription.example pulsar-q4k5898krpqj#tf_example#tf-example-topic#tf-example-subscription#true
```
