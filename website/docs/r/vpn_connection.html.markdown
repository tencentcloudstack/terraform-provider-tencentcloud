---
subcategory: "VPN"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_connection"
sidebar_current: "docs-tencentcloud-resource-vpn_connection"
description: |-
  Provides a resource to create a VPN connection.
---

# tencentcloud_vpn_connection

Provides a resource to create a VPN connection.

## Example Usage

```hcl
resource "tencentcloud_vpn_connection" "foo" {
  name                       = "vpn_connection_test"
  vpc_id                     = "vpc-dk8zmwuf"
  vpn_gateway_id             = "vpngw-8ccsnclt"
  customer_gateway_id        = "cgw-xfqag"
  pre_share_key              = "testt"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "SHA"
  ike_local_identity         = "ADDRESS"
  ike_exchange_mode          = "AGGRESSIVE"
  ike_local_address          = "1.1.1.1"
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = "2.2.2.2"
  ike_dh_group_name          = "GROUP2"
  ike_sa_lifetime_seconds    = 86401
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "SHA1"
  ipsec_sa_lifetime_seconds  = 7200
  ipsec_pfs_dh_group         = "NULL"
  ipsec_sa_lifetime_traffic  = 2570

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["2.2.2.0/26", ]
  }
  tags = {
    test = "testt"
  }
}
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_id` - (Required, String, ForceNew) ID of the customer gateway.
* `name` - (Required, String) Name of the VPN connection. The length of character is limited to 1-60.
* `pre_share_key` - (Required, String) Pre-shared key of the VPN connection.
* `security_group_policy` - (Required, Set) Security group policy of the VPN connection.
* `vpn_gateway_id` - (Required, String, ForceNew) ID of the VPN gateway.
* `enable_health_check` - (Optional, Bool) Whether intra-tunnel health checks are supported.
* `health_check_local_ip` - (Optional, String) Health check the address of this terminal.
* `health_check_remote_ip` - (Optional, String) Health check peer address.
* `ike_dh_group_name` - (Optional, String) DH group name of the IKE operation specification. Valid values: `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`. Default value is `GROUP1`.
* `ike_exchange_mode` - (Optional, String) Exchange mode of the IKE operation specification. Valid values: `AGGRESSIVE`, `MAIN`. Default value is `MAIN`.
* `ike_local_address` - (Optional, String) Local address of IKE operation specification, valid when ike_local_identity is `ADDRESS`, generally the value is `public_ip_address` of the related VPN gateway.
* `ike_local_fqdn_name` - (Optional, String) Local FQDN name of the IKE operation specification.
* `ike_local_identity` - (Optional, String) Local identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.
* `ike_proto_authen_algorithm` - (Optional, String) Proto authenticate algorithm of the IKE operation specification. Valid values: `MD5`, `SHA`, `SHA-256`. Default Value is `MD5`.
* `ike_proto_encry_algorithm` - (Optional, String) Proto encrypt algorithm of the IKE operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-128`, `AES-CBC-256`, `DES-CBC`. Default value is `3DES-CBC`.
* `ike_remote_address` - (Optional, String) Remote address of IKE operation specification, valid when ike_remote_identity is `ADDRESS`, generally the value is `public_ip_address` of the related customer gateway.
* `ike_remote_fqdn_name` - (Optional, String) Remote FQDN name of the IKE operation specification.
* `ike_remote_identity` - (Optional, String) Remote identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.
* `ike_sa_lifetime_seconds` - (Optional, Int) SA lifetime of the IKE operation specification, unit is `second`. The value ranges from 60 to 604800. Default value is 86400 seconds.
* `ike_version` - (Optional, String) Version of the IKE operation specification. Default value is `IKEV1`.
* `ipsec_encrypt_algorithm` - (Optional, String) Encrypt algorithm of the IPSEC operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-128`, `AES-CBC-256`, `DES-CBC`. Default value is `3DES-CBC`.
* `ipsec_integrity_algorithm` - (Optional, String) Integrity algorithm of the IPSEC operation specification. Valid values: `SHA1`, `MD5`, `SHA-256`. Default value is `MD5`.
* `ipsec_pfs_dh_group` - (Optional, String) PFS DH group. Valid value: `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`, `NULL`. Default value is `NULL`.
* `ipsec_sa_lifetime_seconds` - (Optional, Int) SA lifetime of the IPSEC operation specification, unit is second. Valid value ranges: [180~604800]. Default value is 3600 seconds.
* `ipsec_sa_lifetime_traffic` - (Optional, Int) SA lifetime of the IPSEC operation specification, unit is KB. The value should not be less then 2560. Default value is 1843200.
* `tags` - (Optional, Map) A list of tags used to associate different resources.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC. Required if vpn gateway is not in `CCN` type, and doesn't make sense for `CCN` vpn gateway.

The `security_group_policy` object supports the following:

* `local_cidr_block` - (Required, String) Local cidr block.
* `remote_cidr_block` - (Required, Set) Remote cidr block list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the VPN connection.
* `encrypt_proto` - Encrypt proto of the VPN connection.
* `is_ccn_type` - Indicate whether is ccn type. Modification of this field only impacts force new logic of `vpc_id`. If `is_ccn_type` is true, modification of `vpc_id` will be ignored.
* `net_status` - Net status of the VPN connection. Valid value: `AVAILABLE`.
* `route_type` - Route type of the VPN connection.
* `state` - State of the connection. Valid value: `PENDING`, `AVAILABLE`, `DELETING`.
* `vpn_proto` - Vpn proto of the VPN connection.


## Import

VPN connection can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_connection.foo vpnx-nadifg3s
```

