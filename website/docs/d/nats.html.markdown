---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway"
sidebar_current: "docs-tencentcloud-datasource-nats"
description: |-
  The NATs data source lists a number of NATs resource information owned by an TencentCloud account.
---

# tencentcloud_nats

The NATs data source lists a number of NATs resource information owned by an TencentCloud account.

## Example Usage

Basic usage:

```hcl
# Query the NAT gateway by ID
data "tencentcloud_nats" "anat" {
  id = "nat-k6ualnp2"
}

# Query the list of normal NAT gateways
data "tencentcloud_nats" "nat_state" {
  state = 0
}

# Multi conditional query NAT gateway list
data "tencentcloud_nats" "multi_nat" {
  name           = "terraform test"
  vpc_id         = "vpc-ezij4ltv"
  max_concurrent = 3000000
  bandwidth      = 500 
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID for NAT Gateway.
* `name` - (Optional) The name for NAT Gateway.
* `vpc_id` - (Optional) The VPC ID for NAT Gateway.
* `max_concurrent` - (Optional) The upper limit of concurrent connection of NAT gateway, for example: 1000000, 3000000, 10000000. To learn more, please refer to [Virtual Private Cloud Gateway Description](https://intl.cloud.tencent.com/doc/product/215/1682).
* `bandwidth` - (Optional) The maximum public network output bandwidth of the gateway (unit: Mbps), for example: 10, 20, 50, 100, 200, 500, 1000, 2000, 5000. For more information, please refer to [Virtual Private Cloud Gateway Description](https://intl.cloud.tencent.com/doc/product/215/1682).
* `assigned_eip_set` - (Optional) Elastic IP arrays bound to the gateway, For more information on elastic IP, please refer to [Elastic IP](eip.html).
* `state` - (Optional) NAT gateway status, 0: Running, 1: Unavailable, 2: Be in arrears and out of service

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NAT Gateway.
* `name` - The name of the NAT Gateway.
* `max_concurrent` - The upper limit of concurrent connection of the NAT gateway.
* `bandwidth` - The maximum public network output bandwidth of the NAT gateway (unit: Mbps).
* `assigned_eip_set` - Elastic IP arrays bound to the NAT gateway
* `state` - NAT gateway status, 0: Running, 1: Unavailable, 2: Be in arrears and out of service
* `create_time` - The create time of the NAT gateway
