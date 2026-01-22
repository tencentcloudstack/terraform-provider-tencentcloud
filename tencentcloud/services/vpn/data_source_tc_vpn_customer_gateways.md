Use this data source to query detailed information of VPN customer gateways.

Example Usage

Query all customer gateways 

```hcl
data "tencentcloud_vpn_customer_gateways" "example" {}
```

Query customer gateways by filters

```hcl
data "tencentcloud_vpn_customer_gateways" "example" {
  name              = "tf-example"
  id                = "cgw-r1g6c8fr"
  public_ip_address = "1.1.1.1"
  tags = {
    createBy = "Terraform"
  }
}
```
