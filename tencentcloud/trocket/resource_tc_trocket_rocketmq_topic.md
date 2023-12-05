Provides a resource to create a trocket rocketmq_topic

Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxx"
  subnet_id     = "subnet-xxxxx"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_topic" "rocketmq_topic" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  topic       = "test_topic"
  topic_type  = "NORMAL"
  queue_num   = 4
  remark      = "test for terraform"
}
```

Import

trocket rocketmq_topic can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_topic.rocketmq_topic instanceId#topic
```