Provides a resource to create a VPN connection.

Example Usage

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

Import

VPN connection can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_connection.example vpnx-nadifg3s
```
