---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway_snats"
sidebar_current: "docs-tencentcloud-datasource-nat_gateway_snats"
description: |-
  Use this data source to query detailed information of VPN gateways.
---

# tencentcloud_nat_gateway_snats

Use this data source to query detailed information of VPN gateways.

## Example Usage

```hcl
data "tencentcloud_nat_gateway_snats" "snat" {
  nat_gateway_id     = tencentcloud_nat_gateway.my_nat.id
  subnet_id          = tencentcloud_nat_gateway_snat.my_subnet.id
  public_ip_addr     = ["50.29.23.234"]
  description        = "snat demo"
  result_output_file = "./snat.txt"
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required, String) NAT gateway ID.
* `description` - (Optional, String) Description.
* `instance_id` - (Optional, String) Instance ID.
* `public_ip_addr` - (Optional, List: [`String`]) Elastic IP address pool.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) Subnet instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snat_list` - Information list of the nat gateway snat.
  * `create_time` - Create time.
  * `snat_id` - SNAT rule ID.


