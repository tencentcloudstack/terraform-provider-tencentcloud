---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_instance"
sidebar_current: "docs-tencentcloud-resource-trocket_rocketmq_instance"
description: |-
  Provides a resource to create a Trocket rocketmq instance
---

# tencentcloud_trocket_rocketmq_instance

Provides a resource to create a Trocket rocketmq instance

## Example Usage

### Create Basic Instance

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "PRO"
  sku_code      = "pro_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}
```

### Create Enable Public Network Instance

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "PRO"
  sku_code      = "pro_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  enable_public = true
  bandwidth     = 10
  ip_rules {
    ip     = "1.1.1.1"
    allow  = true
    remark = "remark message."
  }

  ip_rules {
    ip     = "2.2.2.2"
    allow  = false
    remark = "remark message."
  }

  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}
```

### Create Instance with Billing and Deployment Params

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rocketmq instance with billing and deployment params
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "PRO"
  sku_code      = "pro_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  pay_mode      = 1
  renew_flag    = 1
  time_span     = 12
  max_topic_num = 1000
  zone_ids      = [100006, 100007]
  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) Instance type. Valid values: `EXPERIMENT` (trial edition), `BASIC` (basic edition), `PRO` (professional edition), `PLATINUM` (platinum edition).
* `name` - (Required, String) Instance (cluster) name, 3-64 characters, can only contain digits, letters, hyphen '-' and underscore '_'.
* `sku_code` - (Required, String) SKU code, obtained from the ProductSKU output of the DescribeProductSKUs interface.
* `subnet_id` - (Required, String) Subnet ID that the instance binds to.
* `vpc_id` - (Required, String) VPC ID that the instance binds to.
* `bandwidth` - (Optional, Int) Public network bandwidth in Mbps, default 0. Must be a positive integer greater than 0 when public network is enabled.
* `enable_public` - (Optional, Bool) Whether to enable public network access, default false. When set to true, `bandwidth` must be set to a positive integer.
* `ip_rules` - (Optional, List) Public network access whitelist. If left empty, all IP access is denied.
* `max_topic_num` - (Optional, Int) Maximum number of topics that can be created. The default/minimum and maximum are obtained from the TopicNumLimit and TopicNumUpperLimit parameters in the ProductSKU output of the DescribeProductSKUs interface.
* `message_retention` - (Optional, Int) Message retention time in hours. The value range and default are obtained from the DefaultRetention/RetentionLowerLimit/RetentionUpperLimit parameters in the ProductSKU output of the DescribeProductSKUs interface.
* `pay_mode` - (Optional, Int) Billing mode. `0`: pay-as-you-go (postpaid), `1`: subscription (prepaid). Default is `0`.
* `remark` - (Optional, String) Remark information.
* `renew_flag` - (Optional, Int) Whether to auto-renew a prepaid instance. `0`: no auto-renewal, `1`: auto-renewal. Default is `0`.
* `tags` - (Optional, Map) Tag list.
* `time_span` - (Optional, Int) Purchase duration of a prepaid instance in months. Value range: 1-60. Default is `1`.
* `zone_ids` - (Optional, List: [`Int`]) List of deployment availability zones, obtained from the ZoneInfo structure returned by the DescribeZones interface.

The `ip_rules` object supports the following:

* `allow` - (Required, Bool) Whether to allow access from this IP.
* `ip` - (Required, String) IP address.
* `remark` - (Required, String) Remark information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `public_end_point` - Public network access address.
* `vpc_end_point` - VPC access address.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateInstance](https://cloud.tencent.com/document/api/1493/97868) - This API is used to create a RocketMQ 5.x cluster.
* Create: [DescribeInstance](https://cloud.tencent.com/document/api/1493/97866) - This API is used to query RocketMQ 5.x cluster information.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeInstance](https://cloud.tencent.com/document/api/1493/97866) - This API is used to query RocketMQ 5.x cluster information.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeInstance](https://cloud.tencent.com/document/api/1493/97866) - This API is used to query RocketMQ 5.x cluster information.
* Update: [ModifyInstance](https://cloud.tencent.com/document/api/1493/97865) - This API is used to modify attributes of a TDMQ RocketMQ 5.x cluster. Only running clusters support this operation.
* Update: [ModifyInstanceEndpoint](https://cloud.tencent.com/document/api/1493/115981) - This API is used to modify access points of a TDMQ RocketMQ 5.x cluster. Make sure that the access points exist before the operation.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteInstance](https://cloud.tencent.com/document/api/1493/97867) - This API is used to delete a TDMQ RocketMQ 5.x cluster. Topics, consumer groups, and roles in use should be deleted first.
* Delete: [DescribeInstance](https://cloud.tencent.com/document/api/1493/97866) - This API is used to query RocketMQ 5.x cluster information.

## Import

Trocket rocketmq instance can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_instance.example rmq-n5qado7m
```

