Provides a resource to create a VPN customer gateway.

Example Usage

```hcl
resource "tencentcloud_vpn_customer_gateway" "example" {
  name              = "tf-example"
  public_ip_address = "1.1.1.1"
  tags = {
    createBy = "Terraform"
  }
}
```

Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_customer_gateway.example cgw-xfqag
```