---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway_flow_monitor"
sidebar_current: "docs-tencentcloud-resource-nat_gateway_flow_monitor"
description: |-
  Provides a resource to config a NAT gateway flow monitor.
---

# tencentcloud_nat_gateway_flow_monitor

Provides a resource to config a NAT gateway flow monitor.

## Example Usage

```hcl
resource "tencentcloud_nat_gateway_flow_monitor" "example" {
  gateway_id = "nat-e6u6axsm"
  enable     = true
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) Whether to enable flow monitor.
* `gateway_id` - (Required, String, ForceNew) ID of Gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bandwidth` - Bandwidth of flow monitor.


## Import

NAT gateway flow monitor can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway_flow_monitor.example nat-e6u6axsm
```

