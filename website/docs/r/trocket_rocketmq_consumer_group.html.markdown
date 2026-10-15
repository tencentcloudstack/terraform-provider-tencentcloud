---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_consumer_group"
sidebar_current: "docs-tencentcloud-resource-trocket_rocketmq_consumer_group"
description: |-
  Provides a resource to create a TROCKET rocketmq consumer group
---

# tencentcloud_trocket_rocketmq_consumer_group

Provides a resource to create a TROCKET rocketmq consumer group

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}

// create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "BASIC"
  sku_code      = "basic_2k"
  remark        = "remark."
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  tags = {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}

// create consumer group
resource "tencentcloud_trocket_rocketmq_consumer_group" "example" {
  instance_id             = tencentcloud_trocket_rocketmq_instance.example.id
  consumer_group          = "tf-example"
  max_retry_times         = 20
  consume_enable          = false
  consume_message_orderly = true
  remark                  = "remark."
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `consume_enable` - (Required, Bool) Whether to enable consumption.
* `consume_message_orderly` - (Required, Bool) `true`: Sequential delivery, `false`: Concurrent delivery.
* `consumer_group` - (Required, String, ForceNew) Name of consumer group.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `max_retry_times` - (Required, Int) Max retry times.
* `remark` - (Optional, String) remark.
* `tags` - (Optional, Map) Tag of consumer group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateConsumerGroup](https://cloud.tencent.com/document/api/1493/97943) - 创建消费组
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeConsumerGroup](https://cloud.tencent.com/document/api/1493/97941) - 查询消费组详情
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [ModifyConsumerGroup](https://cloud.tencent.com/document/api/1493/97940) - 修改消费组属性
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DeleteConsumerGroup](https://cloud.tencent.com/document/api/1493/97942) - 删除消费组

## Import

TROCKET rocketmq consumer group can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_consumer_group.example rmq-1n58qbwg3#tf-example
```

