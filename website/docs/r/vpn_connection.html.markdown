---
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

* `customer_gateway_id` - (Required, ForceNew) ID of the customer gateway.
* `name` - (Required) Name of the VPN connection. The length of character is limited to 1-60.
* `pre_share_key` - (Required) Pre-shared key of the VPN connection.
* `security_group_policy` - (Required) Security group policy of the VPN connection.
* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `vpn_gateway_id` - (Required, ForceNew) ID of the VPN gateway.
* `ike_dh_group_name` - (Optional) DH group name of the IKE operation specification, valid values are `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`. Default value is `GROUP1`.
* `ike_exchange_mode` - (Optional) Exchange mode of the IKE operation specification, valid values are `AGGRESSIVE`, `MAIN`. Default value is `MAIN`.
* `ike_local_address` - (Optional) Local address of IKE operation specification, valid when ike_local_identity is `ADDRESS`, generally the value is public_ip_address of the related VPN gateway.
* `ike_local_fqdn_name` - (Optional) Local FQDN name of the IKE operation specification.
* `ike_local_identity` - (Optional) Local identity way of IKE operation specification, valid values are `ADDRESS`, `FQDN`. Default value is `ADDRESS`.
* `ike_proto_authen_algorithm` - (Optional) Proto authenticate algorithm of the IKE operation specification, valid values are `MD5`, `SHA`. Default Value is `MD5`.
* `ike_proto_encry_algorithm` - (Optional) Proto encrypt algorithm of the IKE operation specification, valid values are `3DES-CBC`, `AES-CBC-128`, `AES-CBC-128`, `AES-CBC-256`, `DES-CBC`. Default value is `3DES-CBC`.
* `ike_remote_address` - (Optional) Remote address of IKE operation specification, valid when ike_remote_identity is `ADDRESS`, generally the value is public_ip_address of the related customer gateway.
* `ike_remote_fqdn_name` - (Optional) Remote FQDN name of the IKE operation specification.
* `ike_remote_identity` - (Optional) Remote identity way of IKE operation specification, valid values are `ADDRESS`, `FQDN`. Default value is `ADDRESS`.
* `ike_sa_lifetime_seconds` - (Optional) SA lifetime of the IKE operation specification, unit is `second`. The value ranges from 60 to 604800. Default value is 86400 seconds.
* `ike_version` - (Optional) Version of the IKE operation specification. Default value is `IKEV1`.
* `ipsec_encrypt_algorithm` - (Optional) Encrypt algorithm of the IPSEC operation specification, valid values are `3DES-CBC`, `AES-CBC-128`, `AES-CBC-128`, `AES-CBC-256`, `DES-CBC`. Default value is `3DES-CBC`.
* `ipsec_integrity_algorithm` - (Optional) Integrity algorithm of the IPSEC operation specification, valid values are `SHA1`, `MD5`. Default value is `MD5`.
* `ipsec_pfs_dh_group` - (Optional) PFS DH group, valid values are `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`, `NULL`. Default value is `NULL`.
* `ipsec_sa_lifetime_seconds` - (Optional) SA lifetime of the IPSEC operation specification, unit is `second`. The value ranges from 180 to 604800. Default value is 3600 seconds.
* `ipsec_sa_lifetime_traffic` - (Optional) SA lifetime of the IPSEC operation specification, unit is `KB`. The value ranges from 2560 to 4294967295. Default value is 1843200.
* `tags` - (Optional) A list of tags used to associate different resources.

The `security_group_policy` object supports the following:

* `local_cidr_block` - (Required) Local cidr block.
* `remote_cidr_block` - (Required) Remote cidr block list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the VPN connection.
* `encrypt_proto` - Encrypt proto of the VPN connection.
* `net_status` - Net status of the VPN connection, values are `AVAILABLE`.
* `route_type` - Route type of the VPN connection.
* `state` - State of the connection, values are `PENDING`, `AVAILABLE`, `DELETING`.
* `vpn_proto` - Vpn proto of the VPN connection.


## Import

VPN connection can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_connection.foo vpnx-nadifg3s
```

