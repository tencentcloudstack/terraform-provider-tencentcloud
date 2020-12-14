---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ha_vips"
sidebar_current: "docs-tencentcloud-datasource-ha_vips"
description: |-
  Use this data source to query detailed information of HA VIPs.
---

# tencentcloud_ha_vips

Use this data source to query detailed information of HA VIPs.

## Example Usage

```hcl
data "tencentcloud_ha_vips" "havips" {
  id         = "havip-kjqwe4ba"
  name       = "test"
  vpc_id     = "vpc-gzea3dd7"
  subnet_id  = "subnet-4d4m4cd4"
  address_ip = "10.0.4.16"
}
```

## Argument Reference

The following arguments are supported:

* `address_ip` - (Optional) EIP of the HA VIP to be queried.
* `id` - (Optional) ID of the HA VIP to be queried.
* `name` - (Optional) Name of the HA VIP. The length of character is limited to 1-60.
* `result_output_file` - (Optional) Used to save results.
* `subnet_id` - (Optional) Subnet id of the HA VIP to be queried.
* `vpc_id` - (Optional) VPC id of the HA VIP to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ha_vip_list` - Information list of the dedicated HA VIPs.
  * `address_ip` - EIP that is associated.
  * `create_time` - Create time of the HA VIP.
  * `id` - ID of the HA VIP.
  * `instance_id` - Instance id that is associated.
  * `name` - Name of the HA VIP.
  * `network_interface_id` - Network interface id that is associated.
  * `state` - State of the HA VIP. Valid values: `AVAILABLE`, `UNBIND`.
  * `subnet_id` - Subnet id.
  * `vip` - Virtual IP address, it must not be occupied and in this VPC network segment. If not set, it will be assigned after resource created automatically.
  * `vpc_id` - VPC id.


