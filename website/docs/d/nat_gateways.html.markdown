---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateways"
sidebar_current: "docs-tencentcloud-datasource-nat_gateways"
description: |-
  Use this data source to query detailed information of NAT gateways.
---

# tencentcloud_nat_gateways

Use this data source to query detailed information of NAT gateways.

## Example Usage

```hcl
data "tencentcloud_nat_gateways" "foo" {
  name   = "main"
  vpc_id = "vpc-xfqag"
  id     = "nat-xfaq1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the nat gateway.
* `name` - (Optional) Name of the nat gateway.
* `result_output_file` - (Optional) Used to save results.
* `vpc_id` - (Optional) ID of the vpc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nats` - Information list of the dedicated tunnels.
  * `assigned_eip_set` - EIP arrays bound to the gateway. The value of at least 1.
  * `bandwidth` - The maximum public network output bandwidth of nat gateway (unit: Mbps), the available values include： 20,50,100,200,500,1000,2000,5000. Default is 100.
  * `create_time` - Create time of the nat gateway.
  * `id` - ID of the nat gateway.
  * `max_concurrent` - The upper limit of concurrent connection of nat gateway, the available values include : 1000000,3000000,10000000, Default is 1000000.
  * `name` - Name of the nat gateway.
  * `state` - State of the nat gateway.
  * `vpc_id` - ID of the vpc.


