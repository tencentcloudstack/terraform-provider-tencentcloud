Provides a resource to create a trocket rocketmq_consumer_group

Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxxx"
  subnet_id     = "subnet-xxxxx"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_consumer_group" "rocketmq_consumer_group" {
  instance_id             = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  consumer_group          = "test_consumer_group"
  max_retry_times         = 20
  consume_enable          = false
  consume_message_orderly = true
  remark                  = "test for terraform"
}
```

Import

trocket rocketmq_consumer_group can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group  instanceId#consumerGroup
```