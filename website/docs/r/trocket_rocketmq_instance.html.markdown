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

~> **NOTE:** It only supports create postpaid rocketmq 5.x instance.

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

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) Instance type. Valid values: `EXPERIMENT`, `BASIC`, `PRO`, `PLATINUM`.
* `name` - (Required, String) Instance name.
* `sku_code` - (Required, String) SKU code. Available specifications are as follows: experiment_500, basic_1k, basic_2k, basic_3k, basic_4k, basic_5k, basic_6k, basic_7k, basic_8k, basic_9k, basic_10k, pro_4k, pro_6k, pro_8k, pro_1w, pro_15k, pro_2w, pro_25k, pro_3w, pro_35k, pro_4w, pro_45k, pro_5w, pro_55k, pro_60k, pro_65k, pro_70k, pro_75k, pro_80k, pro_85k, pro_90k, pro_95k, pro_100k, platinum_1w, platinum_2w, platinum_3w, platinum_4w, platinum_5w, platinum_6w, platinum_7w, platinum_8w, platinum_9w, platinum_10w, platinum_12w, platinum_14w, platinum_16w, platinum_18w, platinum_20w, platinum_25w, platinum_30w, platinum_35w, platinum_40w, platinum_45w, platinum_50w, platinum_60w, platinum_70w, platinum_80w, platinum_90w, platinum_100w.
* `subnet_id` - (Required, String) Subnet id.
* `vpc_id` - (Required, String) VPC id.
* `bandwidth` - (Optional, Int) Public network bandwidth. `bandwidth` must be greater than zero when `enable_public` equal true.
* `enable_public` - (Optional, Bool) Whether to enable the public network. Must set `bandwidth` when `enable_public` equal true.
* `ip_rules` - (Optional, List) Public network access whitelist.
* `message_retention` - (Optional, Int) Message retention time in hours.
* `remark` - (Optional, String) Remark.
* `tags` - (Optional, Map) Tag description list.

The `ip_rules` object supports the following:

* `allow` - (Required, Bool) Whether to allow release or not.
* `ip` - (Required, String) IP.
* `remark` - (Required, String) Remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `public_end_point` - Public network access address.
* `vpc_end_point` - VPC access address.


## Import

Trocket rocketmq instance can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_instance.example rmq-n5qado7m
```

