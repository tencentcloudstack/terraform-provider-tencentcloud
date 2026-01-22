---
subcategory: "VPN Connections(VPN)"
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
resource "tencentcloud_vpn_connection" "example" {
  name                = "tf-example"
  vpc_id              = "vpc-6ccw0s5l"
  vpn_gateway_id      = "vpngw-33p5vnwd"
  customer_gateway_id = "cgw-e503id2z"
  pre_share_key       = "your_pre_share_key"
  route_type          = "StaticRoute"
  negotiation_type    = "flowTrigger"

  # IKE setting
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "SHA"
  ike_local_identity         = "ADDRESS"
  ike_exchange_mode          = "AGGRESSIVE"
  ike_local_address          = "159.75.204.38"
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = "109.244.60.154"
  ike_dh_group_name          = "GROUP2"
  ike_sa_lifetime_seconds    = 86400

  # IPSEC setting
  ipsec_encrypt_algorithm   = "3DES-CBC"
  ipsec_integrity_algorithm = "SHA1"
  ipsec_sa_lifetime_seconds = 14400
  ipsec_pfs_dh_group        = "NULL"
  ipsec_sa_lifetime_traffic = 4096000000

  # health check setting
  enable_health_check    = true
  health_check_local_ip  = "169.254.227.187"
  health_check_remote_ip = "169.254.164.37"
  health_check_config {
    probe_type      = "NQA"
    probe_interval  = 5000
    probe_threshold = 3
    probe_timeout   = 150
  }

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["2.2.2.0/26", ]
  }

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_id` - (Required, String, ForceNew) ID of the customer gateway.
* `name` - (Required, String) Name of the VPN connection. The length of character is limited to 1-60.
* `pre_share_key` - (Required, String) Pre-shared key of the VPN connection.
* `vpn_gateway_id` - (Required, String, ForceNew) ID of the VPN gateway.
* `bgp_config` - (Optional, List, ForceNew) BGP config.
* `dpd_action` - (Optional, String) The action after DPD timeout. Valid values: clear (disconnect) and restart (try again). It is valid when DpdEnable is 1.
* `dpd_enable` - (Optional, Int) Specifies whether to enable DPD. Valid values: 0 (disable) and 1 (enable).
* `dpd_timeout` - (Optional, Int) DPD timeout period.Valid value ranges: [30~60], Default: 30; unit: second. If the request is not responded within this period, the peer end is considered not exists. This parameter is valid when the value of DpdEnable is 1.
* `enable_health_check` - (Optional, Bool) Whether intra-tunnel health checks are supported.
* `health_check_config` - (Optional, List) VPN channel health check configuration.
* `health_check_local_ip` - (Optional, String) Health check the address of this terminal.
* `health_check_remote_ip` - (Optional, String) Health check peer address.
* `ike_dh_group_name` - (Optional, String) DH group name of the IKE operation specification. Valid values: `GROUP1`, `GROUP2`, `GROUP5`, `GROUP14`, `GROUP24`. Default value is `GROUP1`.
* `ike_exchange_mode` - (Optional, String) Exchange mode of the IKE operation specification. Valid values: `AGGRESSIVE`, `MAIN`. Default value is `MAIN`.
* `ike_local_address` - (Optional, String) Local address of IKE operation specification, valid when ike_local_identity is `ADDRESS`, generally the value is `public_ip_address` of the related VPN gateway.
* `ike_local_fqdn_name` - (Optional, String) Local FQDN name of the IKE operation specification.
* `ike_local_identity` - (Optional, String) Local identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.
* `ike_proto_authen_algorithm` - (Optional, String) Proto authenticate algorithm of the IKE operation specification. Valid values: `MD5`, `SHA`, `SHA-256`. Default Value is `MD5`.
* `ike_proto_encry_algorithm` - (Optional, String) Proto encrypt algorithm of the IKE operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-192`, `AES-CBC-256`, `DES-CBC`, `SM4`, `AES128GCM128`, `AES192GCM128`, `AES256GCM128`,`AES128GCM128`, `AES192GCM128`, `AES256GCM128`. Default value is `3DES-CBC`.
* `ike_remote_address` - (Optional, String) Remote address of IKE operation specification, valid when ike_remote_identity is `ADDRESS`, generally the value is `public_ip_address` of the related customer gateway.
* `ike_remote_fqdn_name` - (Optional, String) Remote FQDN name of the IKE operation specification.
* `ike_remote_identity` - (Optional, String) Remote identity way of IKE operation specification. Valid values: `ADDRESS`, `FQDN`. Default value is `ADDRESS`.
* `ike_sa_lifetime_seconds` - (Optional, Int) SA lifetime of the IKE operation specification, unit is `second`. The value ranges from 60 to 604800. Default value is 86400 seconds.
* `ike_version` - (Optional, String) Version of the IKE operation specification, values: `IKEV1`, `IKEV2`. Default value is `IKEV1`.
* `ipsec_encrypt_algorithm` - (Optional, String) Encrypt algorithm of the IPSEC operation specification. Valid values: `3DES-CBC`, `AES-CBC-128`, `AES-CBC-192`, `AES-CBC-256`, `DES-CBC`, `SM4`, `NULL`, `AES128GCM128`, `AES192GCM128`, `AES256GCM128`. Default value is `3DES-CBC`.
* `ipsec_integrity_algorithm` - (Optional, String) Integrity algorithm of the IPSEC operation specification. Valid values: `SHA1`, `MD5`, `SHA-256`. Default value is `MD5`.
* `ipsec_pfs_dh_group` - (Optional, String) PFS DH group. Valid value: `DH-GROUP1`, `DH-GROUP2`, `DH-GROUP5`, `DH-GROUP14`, `DH-GROUP24`, `NULL`. Default value is `NULL`.
* `ipsec_sa_lifetime_seconds` - (Optional, Int) SA lifetime of the IPSEC operation specification, unit is second. Valid value ranges: [180~604800]. Default value is 3600 seconds.
* `ipsec_sa_lifetime_traffic` - (Optional, Int) SA lifetime of the IPSEC operation specification, unit is KB. The value should not be less then 2560. Default value is 1843200.
* `negotiation_type` - (Optional, String) The default negotiation type is `active`. Optional values: `active` (active negotiation), `passive` (passive negotiation), `flowTrigger` (traffic negotiation).
* `route_type` - (Optional, String, ForceNew) Route type of the VPN connection. Valid value: `STATIC`, `StaticRoute`, `Policy`, `Bgp`.
* `security_group_policy` - (Optional, Set) SPD policy group, for example: {"10.0.0.5/24":["172.123.10.5/16"]}, 10.0.0.5/24 is the vpc intranet segment, and 172.123.10.5/16 is the IDC network segment. Users specify which network segments in the VPC can communicate with which network segments in your IDC.
* `tags` - (Optional, Map) A list of tags used to associate different resources.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC. Required if vpn gateway is not in `CCN` type, and doesn't make sense for `CCN` vpn gateway.

The `bgp_config` object supports the following:

* `local_bgp_ip` - (Required, String) Cloud BGP address. It must be allocated from within the BGP tunnel network segment.
* `remote_bgp_ip` - (Required, String) User side BGP address. It must be allocated from within the BGP tunnel network segment.
* `tunnel_cidr` - (Required, String) BGP tunnel segment.

The `health_check_config` object supports the following:

* `probe_interval` - (Optional, Int) Detection interval, Tencent Cloud's interval between two health checks, range [1000-5000], Unit: ms.
* `probe_threshold` - (Optional, Int) Detection times, perform route switching after N consecutive health check failures, range [3-8], Unit: times.
* `probe_timeout` - (Optional, Int) Detection timeout, range [10-5000], Unit: ms.
* `probe_type` - (Optional, String) Detection mode, default is `NQA`, cannot be modified.

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
* `state` - State of the connection. Valid value: `PENDING`, `AVAILABLE`, `DELETING`.
* `vpn_proto` - Vpn proto of the VPN connection.


## Import

VPN connection can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_connection.example vpnx-nadifg3s
```

