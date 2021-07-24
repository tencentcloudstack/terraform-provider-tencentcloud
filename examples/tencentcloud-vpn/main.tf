
resource "tencentcloud_vpn_customer_gateway" "example" {
  name              = "example"
  public_ip_address = "3.3.3.3"

  tags = {
    test = "example"
  }
}

data "tencentcloud_vpc_instances" "example" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "example" {
  name      = "example"
  vpc_id    = data.tencentcloud_vpc_instances.example.instance_list.0.vpc_id
  bandwidth = 5
  zone      = var.availability_zone

  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpn_connection" "example" {
  name                       = "example"
  vpc_id                     = data.tencentcloud_vpc_instances.example.instance_list.0.vpc_id
  vpn_gateway_id             = tencentcloud_vpn_gateway.example.id
  customer_gateway_id        = tencentcloud_vpn_customer_gateway.example.id
  pre_share_key              = "test"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "MD5"
  ike_local_identity         = "ADDRESS"
  ike_local_address          = tencentcloud_vpn_gateway.example.public_ip_address
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = tencentcloud_vpn_customer_gateway.example.public_ip_address
  ike_dh_group_name          = "GROUP1"
  ike_sa_lifetime_seconds    = 86400
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "MD5"
  ipsec_sa_lifetime_seconds  = 3600
  ipsec_pfs_dh_group         = "DH-GROUP1"
  ipsec_sa_lifetime_traffic  = 2560

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["3.3.3.0/32", ]
  }
  tags = {
    test = "test"
  }
}

data "tencentcloud_vpn_customer_gateways" "example" {
  id = tencentcloud_vpn_customer_gateway.example.id
}

data "tencentcloud_vpn_gateways" "example" {
  id = tencentcloud_vpn_gateway.example.id
}

data "tencentcloud_vpn_connections" "example" {
  id = tencentcloud_vpn_connection.example.id
}

# The example below shows how to create a vpn gateway in ccn type if it is needed. Then could be used when creating
# vpn tunnel in the usual way.
resource tencentcloud_vpn_gateway ccn_vpngw_example {
  name      = "ccn-vpngw-example"
  vpc_id    = data.tencentcloud_vpc_instances.example.instance_list.0.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "CCN"

  tags = {
    test = "ccn-vpngw-example"
  }
}

resource "tencentcloud_vpn_gateway_route" "example" {
    vpn_gateway_id = tencentcloud_vpn_gateway.example.id
    destination_cidr_block = "10.0.0.0/16"
    instance_type = "VPNCONN"
    instance_id = "vpnx-5b5dmao3"
    priority = "100"
    status = "DISABLE"
}

data "tencentcloud_vpn_gateway_routes" "example" {
  vpn_gateway_id = tencentcloud_vpn_gateway.example.id
}