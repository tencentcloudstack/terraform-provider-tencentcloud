Provides a resource to create a rocketmq 5.x instance

~> **NOTE:** It only support create postpaid rocketmq 5.x instance.

Example Usage

Basic Instance
```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name = "rocketmq-instance"
  sku_code = "experiment_500"
  remark = "remark"
  vpc_id = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x"
  }
}
```

Enable Public Instance
```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance_public" {
  instance_type = "EXPERIMENT"
  name = "rocketmq-enable-public-instance"
  sku_code = "experiment_500"
  remark = "remark"
  vpc_id = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x"
  }
  enable_public = true
  bandwidth = 1
}
```

Import

trocket rocketmq_instance can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_instance.rocketmq_instance rocketmq_instance_id
```