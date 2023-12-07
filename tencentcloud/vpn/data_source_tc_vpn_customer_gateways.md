Use this data source to query detailed information of VPN customer gateways.

Example Usage

```hcl
data "tencentcloud_customer_gateways" "foo" {
  name              = "main"
  id                = "cgw-xfqag"
  public_ip_address = "1.1.1.1"
  tags = {
    test = "tf"
  }
}
```