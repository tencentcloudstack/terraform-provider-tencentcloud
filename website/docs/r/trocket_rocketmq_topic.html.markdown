---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_topic"
sidebar_current: "docs-tencentcloud-resource-trocket_rocketmq_topic"
description: |-
  Provides a resource to create a TROCKET rocketmq topic
---

# tencentcloud_trocket_rocketmq_topic

Provides a resource to create a TROCKET rocketmq topic

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

// create topic
resource "tencentcloud_trocket_rocketmq_topic" "example" {
  instance_id = tencentcloud_trocket_rocketmq_instance.example.id
  topic       = "tf-example"
  topic_type  = "NORMAL"
  queue_num   = 4
  remark      = "remark."
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance Id.
* `queue_num` - (Required, Int) Number of queue. Must be greater than or equal to 3.
* `topic_type` - (Required, String, ForceNew) Topic type. `UNSPECIFIED`: not specified, `NORMAL`: normal message, `FIFO`: sequential message, `DELAY`: delayed message.
* `topic` - (Required, String, ForceNew) topic.
* `remark` - (Optional, String) remark.
* `tags` - (Optional, Map) Tag of topic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateTopic](https://www.tencentcloud.com/search/?wd=CreateTopic&type=document&lang=en)
  - [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read:
  - [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
  - [DescribeTopicList](https://www.tencentcloud.com/document/api/1113/74624)
* Update:
  - [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
  - [ModifyTopic](https://www.tencentcloud.com/document/api/1113/74622)
* Delete: [DeleteTopic](https://www.tencentcloud.com/document/api/1113/74627)

## Import

TROCKET rocketmq topic can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_topic.example rmq-1zj5vokgn#tf-example
```

