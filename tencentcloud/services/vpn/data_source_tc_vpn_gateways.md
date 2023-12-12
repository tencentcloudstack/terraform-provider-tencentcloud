Use this data source to query detailed information of VPN gateways.

Example Usage

```hcl
data "tencentcloud_vpn_gateways" "foo" {
  name              = "main"
  id                = "vpngw-8ccsnclt"
  public_ip_address = "1.1.1.1"
  zone              = "ap-guangzhou-3"
  vpc_id            = "vpc-dk8zmwuf"
  tags = {
    test = "tf"
  }
}
```