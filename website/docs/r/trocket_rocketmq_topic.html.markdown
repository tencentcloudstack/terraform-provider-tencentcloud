---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_topic"
sidebar_current: "docs-tencentcloud-resource-trocket_rocketmq_topic"
description: |-
  Provides a resource to create a trocket rocketmq_topic
---

# tencentcloud_trocket_rocketmq_topic

Provides a resource to create a trocket rocketmq_topic

## Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-xxxxx"
  subnet_id     = "subnet-xxxxx"
  tags = {
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

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance Id.
* `queue_num` - (Required, Int) Number of queue. Must be greater than or equal to 3.
* `topic_type` - (Required, String, ForceNew) Topic type. `UNSPECIFIED`: not specified, `NORMAL`: normal message, `FIFO`: sequential message, `DELAY`: delayed message.
* `topic` - (Required, String, ForceNew) topic.
* `remark` - (Optional, String) remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

trocket rocketmq_topic can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_topic.rocketmq_topic instanceId#topic
```

