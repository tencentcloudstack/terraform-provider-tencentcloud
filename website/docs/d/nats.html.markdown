---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nats"
sidebar_current: "docs-tencentcloud-datasource-nats"
description: |-
  The NATs data source lists a number of NATs resource information owned by an TencentCloud account.
---

# tencentcloud_nats

The NATs data source lists a number of NATs resource information owned by an TencentCloud account.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_nat_gateways.

## Example Usage

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

* `bandwidth` - (Optional) The maximum public network output bandwidth of the gateway (unit: Mbps), for example: 10, 20, 50, 100, 200, 500, 1000, 2000, 5000.
* `id` - (Optional) The ID for NAT Gateway.
* `max_concurrent` - (Optional) The upper limit of concurrent connection of NAT gateway, for example: 1000000, 3000000, 10000000.
* `name` - (Optional) The name for NAT Gateway.
* `state` - (Optional) NAT gateway status. Valid values: 0, 1, 2. 0: Running, 1: Unavailable, 2: Be in arrears and out of service.
* `vpc_id` - (Optional) The VPC ID for NAT Gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nats` - Information list of the dedicated tunnels.
  * `assigned_eip_set` - Elastic IP arrays bound to the gateway.
  * `bandwidth` - The maximum public network output bandwidth of the gateway (unit: Mbps), for example: 10, 20, 50, 100, 200, 500, 1000, 2000, 5000.
  * `create_time` - The create time of the NAT gateway.
  * `id` - The ID for NAT Gateway.
  * `max_concurrent` - The upper limit of concurrent connection of NAT gateway, for example: 1000000, 3000000, 10000000.
  * `name` - The name for NAT Gateway.
  * `state` - NAT gateway status, 0: Running, 1: Unavailable, 2: Be in arrears and out of service.
  * `vpc_id` - The VPC ID for NAT Gateway.


