---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_instance"
sidebar_current: "docs-tencentcloud-resource-trocket_rocketmq_instance"
description: |-
  Provides a resource to create a rocketmq 5.x instance
---

# tencentcloud_trocket_rocketmq_instance

Provides a resource to create a rocketmq 5.x instance

~> **NOTE:** It only support create postpaid rocketmq 5.x instance.

## Example Usage

### Basic Instance

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "rocketmq-instance"
  sku_code      = "experiment_500"
  remark        = "remark"
  vpc_id        = "vpc-xxxxxx"
  subnet_id     = "subnet-xxxxxx"
  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}
```

### Enable Public Instance

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance_public" {
  instance_type = "EXPERIMENT"
  name          = "rocketmq-enable-public-instance"
  sku_code      = "experiment_500"
  remark        = "remark"
  vpc_id        = "vpc-xxxxxx"
  subnet_id     = "subnet-xxxxxx"
  tags = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
  enable_public = true
  bandwidth     = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) Instance type. Valid values: `EXPERIMENT`, `BASIC`, `PRO`, `PLATINUM`.
* `name` - (Required, String) Instance name.
* `sku_code` - (Required, String) SKU code. Available specifications are as follows: experiment_500, basic_1k, basic_2k, basic_4k, basic_6k.
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

trocket rocketmq_instance can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_instance.rocketmq_instance rocketmq_instance_id
```

