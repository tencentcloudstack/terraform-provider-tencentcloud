Provide a resource to create a TDMQ topic.

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
  environ_id        = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id        = tencentcloud_tdmq_instance.example.id
  topic_name        = "tf-example-topic"
  partitions        = 6
  pulsar_topic_type = 3
  remark            = "remark."
}
```

Import

Tdmq Topic can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_topic.test topic_id
```