resource "tencentcloud_vpn_customer_gateway" "example" {
  name              = var.vpn_cgw
  public_ip_address = "6.6.6.6"
}

resource "tencentcloud_vpn_gateway" "example" {
  bandwidth          = 5
  charge_type        = "POSTPAID_BY_HOUR"
  name               = var.vpn_gw
  prepaid_renew_flag = "NOTIFY_AND_AUTO_RENEW"
  type               = "IPSEC"
  zone               = "ap-guangzhou-6"
}

resource "tencentcloud_vpn_connection" "example" {
  customer_gateway_id        = resource.tencentcloud_vpn_customer_gateway.example.id
  ike_dh_group_name          = "GROUP1"
  ike_exchange_mode          = "MAIN"
  ike_local_address          = "43.132.156.41"
  ike_local_identity         = "ADDRESS"
  ike_proto_authen_algorithm = "MD5"
  ike_proto_encry_algorithm  = "AES-CBC-128"
  ike_remote_address         = "6.6.6.6"
  ike_sa_lifetime_seconds    = 86400
  ike_version                = "IKEV1"
  ipsec_encrypt_algorithm    = "AES-CBC-128"
  ipsec_integrity_algorithm  = "MD5"
  ipsec_pfs_dh_group         = "NULL"
  ipsec_sa_lifetime_seconds  = 3600
  ipsec_sa_lifetime_traffic  = 1843200
  name                       = var.vpn_conn
  pre_share_key              = "test"
  tags                       = {}
  vpc_id                     = "vpc-4owdpnwr"
  vpn_gateway_id             =  resource.tencentcloud_vpn_gateway.example.id

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = [
      "3.3.3.0/26",
    ]
  }
}
