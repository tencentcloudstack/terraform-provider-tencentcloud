Use this data source to query detailed information of VPN connections.

Example Usage

Query all vpn connections

```hcl
data "tencentcloud_vpn_connections" "example" {}
```

Query vpn connections by filters

```hcl
data "tencentcloud_vpn_connections" "example" {
  name                = "tf-example"
  id                  = "vpnx-fq4e4364"
  vpn_gateway_id      = "vpngw-8ccsnclt"
  vpc_id              = "vpc-6ccw0s5l"
  customer_gateway_id = "cgw-r1g6c8fr"
  tags = {
    createBy = "Terraform"
  }
}
```
