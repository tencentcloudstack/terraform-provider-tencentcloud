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

* Create: [CreateTopic](https://cloud.tencent.com/document/api/1493/97947)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeTopicList](https://cloud.tencent.com/document/api/1493/96030) - This API is used to search the topic list. Filter parameter usage instructions are as follows:.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Update: [ModifyTopic](https://cloud.tencent.com/document/api/1493/97944)
* Delete: [DeleteTopic](https://cloud.tencent.com/document/api/1493/97946) - This API is used to delete a topic. After deletion, all configurations and related data of the topic will be cleared and cannot be restored.

## Import

TROCKET rocketmq topic can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_topic.example rmq-1zj5vokgn#tf-example
```

