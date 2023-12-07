Provides a resource to create a VPN connection.

Example Usage

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

Import

VPN connection can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_connection.foo vpnx-nadifg3s
```