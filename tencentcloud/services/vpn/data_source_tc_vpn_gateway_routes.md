Use this data source to query detailed information of VPN gateways.

Example Usage

```hcl
data "tencentcloud_vpn_gateways" "foo" {
  vpn_gateway_id              = "main"
  destination_cidr_block                = "vpngw-8ccsnclt"
  instance_type = "1.1.1.1"
  instance_id              = "ap-guangzhou-3"
  tags = {
    test = "tf"
  }
}
```