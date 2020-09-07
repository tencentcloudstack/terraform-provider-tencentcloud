---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcx_instances"
sidebar_current: "docs-tencentcloud-datasource-dcx_instances"
description: |-
  Use this data source to query detailed information of dedicated tunnels instances.
---

# tencentcloud_dcx_instances

Use this data source to query detailed information of dedicated tunnels instances.

## Example Usage

```hcl
data "tencentcloud_dcx_instances" "name_select" {
  name = "main"
}

data "tencentcloud_dcx_instances" "id" {
  dcx_id = "dcx-3ikuw30k"
}
```

## Argument Reference

The following arguments are supported:

* `dcx_id` - (Optional) ID of the dedicated tunnels to be queried.
* `name` - (Optional) Name of the dedicated tunnels to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Information list of the dedicated tunnels.
  * `bandwidth` - Bandwidth of the DC.
  * `bgp_asn` - BGP ASN of the user.
  * `bgp_auth_key` - BGP key of the user.
  * `create_time` - Creation time of resource.
  * `customer_address` - Interconnect IP of the DC within client.
  * `dc_id` - ID of the DC.
  * `dcg_id` - ID of the DC Gateway. Currently only new in the console.
  * `dcx_id` - ID of the dedicated tunnel.
  * `name` - Name of the dedicated tunnel.
  * `network_region` - The region of the dedicated tunnel.
  * `network_type` - Type of the network, and available values include VPC, BMVPC and CCN. The default value is VPC.
  * `route_filter_prefixes` - Static route, the network address of the user IDC.
  * `route_type` - Type of the route, and available values include BGP and STATIC. The default value is BGP.
  * `state` - State of the dedicated tunnels, and available values include PENDING, ALLOCATING, ALLOCATED, ALTERING, DELETING, DELETED, COMFIRMING and REJECTED.
  * `tencent_address` - Interconnect IP of the DC within Tencent.
  * `vlan` - Vlan of the dedicated tunnels, and the range of values is [0-3000]. '0' means that only one tunnel can be created for the physical connect.
  * `vpc_id` - ID of the VPC or BMVPC.


