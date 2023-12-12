Provides a resource to create a VPN customer gateway.

Example Usage

```hcl
resource "tencentcloud_vpn_customer_gateway" "foo" {
  name              = "test_vpn_customer_gateway"
  public_ip_address = "1.1.1.1"

  tags = {
    tag = "test"
  }
}
```

Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_customer_gateway.foo cgw-xfqag
```