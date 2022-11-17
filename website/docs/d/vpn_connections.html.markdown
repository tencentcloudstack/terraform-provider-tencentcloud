---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_connections"
sidebar_current: "docs-tencentcloud-datasource-vpn_connections"
description: |-
  Use this data source to query detailed information of VPN connections.
---

# tencentcloud_vpn_connections

Use this data source to query detailed information of VPN connections.

## Example Usage

```hcl
data "tencentcloud_vpn_connections" "foo" {
  name                = "main"
  id                  = "vpnx-xfqag"
  vpn_gateway_id      = "vpngw-8ccsnclt"
  vpc_id              = "cgw-xfqag"
  customer_gateway_id = ""
  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_id` - (Optional, String) Customer gateway ID of the VPN connection.
* `id` - (Optional, String) ID of the VPN connection.
* `name` - (Optional, String) Name of the VPN connection. The length of character is limited to 1-60.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tags of the VPN connection to be queried.
* `vpc_id` - (Optional, String) ID of the VPC.
* `vpn_gateway_id` - (Optional, String) VPN gateway ID of the VPN connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `connection_list` - Information list of the dedicated connections.
  * `create_time` - Create time of the VPN connection.
  * `customer_gateway_id` - ID of the customer gateway.
  * `encrypt_proto` - Encrypt proto of the VPN connection.
  * `id` - ID of the VPN connection.
  * `ike_dh_group_name` - DH group name of the IKE operation specification.
  * `ike_exchange_mode` - Exchange mode of the IKE operation specification.
  * `ike_local_address` - Local address of the IKE operation specification.
  * `ike_local_fqdn_name` - Local FQDN name of the IKE operation specification.
  * `ike_local_identity` - Local identity of the IKE operation specification.
  * `ike_proto_authen_algorithm` - Proto authenticate algorithm of the IKE operation specification.
  * `ike_proto_encry_algorithm` - Proto encrypt algorithm of the IKE operation specification.
  * `ike_remote_address` - Remote address of the IKE operation specification.
  * `ike_remote_fqdn_name` - Remote FQDN name of the IKE operation specification.
  * `ike_remote_identity` - Remote identity of the IKE operation specification.
  * `ike_sa_lifetime_seconds` - SA lifetime of the IKE operation specification, unit is `second`.
  * `ike_version` - Version of the IKE operation specification.
  * `ipsec_encrypt_algorithm` - Encrypt algorithm of the IPSEC operation specification.
  * `ipsec_integrity_algorithm` - Integrity algorithm of the IPSEC operation specification.
  * `ipsec_pfs_dh_group` - PFS DH group name of the IPSEC operation specification.
  * `ipsec_sa_lifetime_seconds` - SA lifetime of the IPSEC operation specification, unit is `second`.
  * `ipsec_sa_lifetime_traffic` - SA lifetime traffic of the IPSEC operation specification, unit is `KB`.
  * `name` - Name of the VPN connection.
  * `net_status` - Net status of the VPN connection.
  * `pre_share_key` - Pre-shared key of the VPN connection.
  * `route_type` - Route type of the VPN connection.
  * `security_group_policy` - Security group policy of the VPN connection.
    * `local_cidr_block` - Local cidr block.
    * `remote_cidr_block` - Remote cidr block list.
  * `state` - State of the VPN connection.
  * `tags` - A list of tags used to associate different resources.
  * `vpc_id` - ID of the VPC.
  * `vpn_gateway_id` - ID of the VPN gateway.
  * `vpn_proto` - Vpn proto of the VPN connection.


